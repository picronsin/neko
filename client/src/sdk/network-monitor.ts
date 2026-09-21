export type NetworkQuality = 'unknown' | 'good' | 'fair' | 'poor'
export type NetworkPath = 'unknown' | 'direct' | 'relay'
export type NetworkProtocol = 'unknown' | 'udp' | 'tcp'

export interface NetworkQualitySample {
  quality: NetworkQuality
  rtt: number | null
  packetLoss: number
  packetsReceived: number
  packetsLost: number
  path: NetworkPath
  protocol: NetworkProtocol
}

export interface NetworkStatsPeer {
  getStats: () => Promise<RTCStatsReport>
}

export interface NetworkQualityMonitorOptions {
  intervalMs?: number
  onSample: (sample: NetworkQualitySample) => void
}

/**
 * Samples inbound video and candidate-pair statistics without depending on
 * Vuex or UI components. A missing peer or transient getStats failure is
 * intentionally ignored; the next sample can recover the display.
 */
export class NetworkQualityMonitor {
  private readonly intervalMs: number
  private readonly onSample: (sample: NetworkQualitySample) => void
  private peer?: NetworkStatsPeer
  private timer?: number
  private previous?: { packetsReceived: number; packetsLost: number }
  private generation = 0
  private sampling = false

  constructor(options: NetworkQualityMonitorOptions) {
    this.intervalMs = options.intervalMs ?? 5000
    this.onSample = options.onSample
  }

  start(peer: NetworkStatsPeer) {
    this.stop()
    this.peer = peer
    this.previous = undefined
    void this.sample()
    this.timer = window.setInterval(() => void this.sample(), this.intervalMs)
  }

  stop() {
    this.generation++
    this.sampling = false
    if (this.timer !== undefined) {
      window.clearInterval(this.timer)
      this.timer = undefined
    }
    this.peer = undefined
    this.previous = undefined
  }

  private async sample() {
    if (!this.peer || this.sampling) {
      return
    }

    const generation = this.generation
    this.sampling = true
    try {
      const stats = await this.peer.getStats()
      if (generation !== this.generation) {
        return
      }
      let packetsReceived = 0
      let packetsLost = 0
      let rtt: number | null = null

      stats.forEach((stat: any) => {
        if (stat.type === 'inbound-rtp' && (stat.kind === 'video' || stat.mediaType === 'video')) {
          packetsReceived += Number(stat.packetsReceived || 0)
          packetsLost += Number(stat.packetsLost || 0)
        }

        if (
          stat.type === 'candidate-pair' &&
          (stat.state === 'succeeded' || stat.nominated === true) &&
          typeof stat.currentRoundTripTime === 'number'
        ) {
          rtt = stat.currentRoundTripTime * 1000
        }
      })

      const selectedPath = selectedCandidatePath(Array.from(stats.values()))
      if (selectedPath.rtt !== null) {
        rtt = selectedPath.rtt
      }

      const previous = this.previous
      this.previous = { packetsReceived, packetsLost }
      const receivedDelta = previous ? Math.max(0, packetsReceived - previous.packetsReceived) : packetsReceived
      const lostDelta = previous ? Math.max(0, packetsLost - previous.packetsLost) : packetsLost
      const totalPackets = receivedDelta + lostDelta
      const packetLoss = totalPackets > 0 ? lostDelta / totalPackets : 0
      const quality = classifyNetworkQuality(rtt, packetLoss, totalPackets > 0 || rtt !== null)

      this.onSample({
        quality,
        rtt: rtt === null ? null : Math.round(rtt),
        packetLoss,
        packetsReceived,
        packetsLost,
        path: selectedPath.path,
        protocol: selectedPath.protocol,
      })
    } catch {
      // getStats is best effort; a temporary failure must not affect the media session.
    } finally {
      if (generation === this.generation) {
        this.sampling = false
      }
    }
  }
}

/**
 * Reads the browser-selected ICE candidate pair without exposing candidate
 * addresses. A relay candidate means media traverses TURN; every other
 * selected pair is a browser-to-server direct path.
 */
export function selectedCandidatePath(stats: Iterable<any>): {
  path: NetworkPath
  protocol: NetworkProtocol
  rtt: number | null
} {
  const values = Array.from(stats)
  const byID = new Map(values.filter((value) => typeof value?.id === 'string').map((value) => [value.id, value]))
  const transport = values.find((value) => value?.type === 'transport' && value.selectedCandidatePairId)
  const succeededPairs = values
    .filter((value) => value?.type === 'candidate-pair' && value.state === 'succeeded')
    .sort(
      (left, right) =>
        Number(right.bytesSent || 0) +
        Number(right.bytesReceived || 0) -
        Number(left.bytesSent || 0) -
        Number(left.bytesReceived || 0),
    )
  const pair =
    byID.get(transport?.selectedCandidatePairId) ||
    succeededPairs.find((value) => value.selected === true || value.nominated === true) ||
    // Some Chromium versions expose neither selected nor nominated. The
    // successful pair carrying the most traffic is the active media path.
    succeededPairs[0]

  if (!pair) {
    return { path: 'unknown', protocol: 'unknown', rtt: null }
  }

  const local = byID.get(pair.localCandidateId)
  const remote = byID.get(pair.remoteCandidateId)
  const protocol = String(pair.protocol || local?.protocol || remote?.protocol || '').toLowerCase()
  return {
    path: local?.candidateType === 'relay' || remote?.candidateType === 'relay' ? 'relay' : 'direct',
    protocol: protocol === 'udp' || protocol === 'tcp' ? protocol : 'unknown',
    rtt: typeof pair.currentRoundTripTime === 'number' ? Math.round(pair.currentRoundTripTime * 1000) : null,
  }
}

export function classifyNetworkQuality(rtt: number | null, packetLoss: number, hasStats: boolean): NetworkQuality {
  if (!hasStats) {
    return 'unknown'
  }

  if ((rtt !== null && rtt > 350) || packetLoss > 0.08) {
    return 'poor'
  }

  if ((rtt !== null && rtt > 180) || packetLoss > 0.03) {
    return 'fair'
  }

  return 'good'
}
