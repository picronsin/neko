export interface VoiceMember {
  id: string
  name: string
  muted: boolean
}
interface Signal {
  from: string
  to: string
  type: RTCSdpType
  sdp: string
}
interface VoiceHTTP {
  get<T>(path: string): Promise<{ data: T }>
  post<T>(path: string, body?: unknown): Promise<{ data: T }>
  delete(path: string): Promise<unknown>
}
interface VoicePeer {
  connection: RTCPeerConnection
  audio: HTMLAudioElement
}

/** Independent room call: never attaches tracks to the remote desktop peer. */
export class VoiceCall {
  private stream?: MediaStream
  private peers = new Map<string, VoicePeer>()
  private id = ''
  private servers: RTCIceServer[] = []
  private timer?: ReturnType<typeof setTimeout>
  private generation = 0
  muted = false

  constructor(
    private http: VoiceHTTP,
    private members: (members: VoiceMember[]) => void,
    private error: (error: Error) => void,
  ) {}

  async join() {
    const generation = ++this.generation
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
        video: false,
      })
      if (generation !== this.generation) {
        stream.getTracks().forEach((track) => track.stop())
        return
      }
      this.stream = stream
      this.muted = false
      const { data } = await this.http.post<{ id: string; servers: RTCIceServer[] }>('/api/chat/voice')
      if (generation !== this.generation) {
        await this.http.delete('/api/chat/voice')
        return
      }
      this.id = data.id
      this.servers = data.servers || []
      void this.poll(generation)
    } catch (error) {
      await this.leave()
      throw error
    }
  }

  setMuted(muted: boolean) {
    this.muted = muted
    this.stream?.getAudioTracks().forEach((track) => {
      track.enabled = !muted
    })
  }

  private peer(id: string): VoicePeer {
    const existing = this.peers.get(id)
    if (existing) return existing
    const connection = new RTCPeerConnection({ iceServers: this.servers })
    const audio = new Audio()
    audio.autoplay = true
    connection.ontrack = ({ track }) => {
      audio.srcObject = new MediaStream([track])
      void audio.play().catch(() => this.fail(new Error('Allow audio playback to hear the voice call')))
    }
    connection.onconnectionstatechange = () => {
      if (connection.connectionState === 'failed')
        this.fail(new Error('Voice connection failed; check TURN configuration and rejoin'))
    }
    this.stream?.getTracks().forEach((track) => connection.addTrack(track, this.stream!))
    const peer = { connection, audio }
    this.peers.set(id, peer)
    return peer
  }

  private async describe(id: string, type: 'offer' | 'answer', generation: number) {
    const { connection } = this.peer(id)
    await connection.setLocalDescription(
      type === 'offer' ? await connection.createOffer() : await connection.createAnswer(),
    )
    // Send one complete description so no candidate ordering/race is exposed
    // to the REST signaling transport.
    if (connection.iceGatheringState !== 'complete') {
      await new Promise<void>((resolve, reject) => {
        const finish = () => {
          if (connection.iceGatheringState !== 'complete' && connection.signalingState !== 'closed') return
          clearTimeout(timer)
          connection.removeEventListener('icegatheringstatechange', finish)
          connection.removeEventListener('signalingstatechange', finish)
          resolve()
        }
        const timer = setTimeout(() => {
          connection.removeEventListener('icegatheringstatechange', finish)
          connection.removeEventListener('signalingstatechange', finish)
          reject(new Error('Voice ICE gathering timed out'))
        }, 10000)
        connection.addEventListener('icegatheringstatechange', finish)
        connection.addEventListener('signalingstatechange', finish)
        finish()
      })
    }
    if (generation !== this.generation || connection.signalingState === 'closed') return
    await this.http.post('/api/chat/voice/signal', { to: id, type, sdp: connection.localDescription!.sdp })
  }

  private async poll(generation: number) {
    try {
      const { data } = await this.http.get<{ members: VoiceMember[]; signals: Signal[] }>(
        `/api/chat/voice?muted=${this.muted}`,
      )
      if (generation !== this.generation) return
      this.members(data.members)
      const ids = new Set(data.members.map((member) => member.id))
      for (const id of this.peers.keys()) if (!ids.has(id)) this.closePeer(id)
      for (const signal of data.signals || []) {
        if (!ids.has(signal.from)) continue
        const { connection } = this.peer(signal.from)
        await connection.setRemoteDescription({ type: signal.type, sdp: signal.sdp })
        if (signal.type === 'offer')
          void this.describe(signal.from, 'answer', generation).catch((error) => {
            if (generation === this.generation) this.fail(error)
          })
      }
      for (const member of data.members) {
        // Exactly one offerer per pair avoids simultaneous-offer glare.
        if (member.id > this.id && !this.peers.has(member.id)) {
          void this.describe(member.id, 'offer', generation).catch((error) => {
            if (generation === this.generation) this.fail(error)
          })
        }
      }
      if (generation === this.generation) this.timer = setTimeout(() => void this.poll(generation), 1000)
    } catch (error) {
      if (generation === this.generation) this.fail(error as Error)
    }
  }

  private fail(error: Error) {
    void this.leave()
    this.error(error)
  }
  private closePeer(id: string) {
    const peer = this.peers.get(id)
    if (!peer) return
    peer.connection.onconnectionstatechange = null
    peer.connection.ontrack = null
    peer.connection.close()
    peer.audio.pause()
    peer.audio.srcObject = null
    this.peers.delete(id)
  }

  async leave() {
    ++this.generation
    clearTimeout(this.timer)
    this.stream?.getTracks().forEach((track) => track.stop())
    this.stream = undefined
    for (const id of this.peers.keys()) this.closePeer(id)
    this.members([])
    if (this.id) {
      this.id = ''
      await this.http.delete('/api/chat/voice').catch(() => {})
    }
  }
}
