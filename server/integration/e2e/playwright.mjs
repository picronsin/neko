import { chromium } from 'playwright'
import { mkdir, writeFile } from 'node:fs/promises'
import process from 'node:process'

const baseURL = process.env.NEKO_E2E_BASE_URL || 'http://127.0.0.1:8080'
const username = process.env.NEKO_E2E_USERNAME || 'e2e-browser'
const password = process.env.NEKO_E2E_PASSWORD
const timeout = Number(process.env.NEKO_E2E_TIMEOUT_MS || 45_000)
const artifactDir = process.env.NEKO_E2E_ARTIFACT_DIR || '/tmp/neko-e2e-artifacts'
const outputPath = process.env.NEKO_E2E_OUTPUT || ''
const viewportWidth = Number(process.env.NEKO_E2E_VIEWPORT_WIDTH || 1280)
const viewportHeight = Number(process.env.NEKO_E2E_VIEWPORT_HEIGHT || 720)

if (!password) {
  throw new Error('NEKO_E2E_PASSWORD must be set')
}

const events = []
let malformedFrames = 0
let debugFrameCount = 0
const startedAt = performance.now()
let browser
let page

function parseFrame(data) {
  if (process.env.NEKO_E2E_DEBUG_WS === '1' && debugFrameCount < 8) {
    debugFrameCount += 1
    const keys = data && typeof data === 'object' ? Object.keys(data).join(',') : ''
    const valueTypes =
      data && typeof data === 'object'
        ? Object.entries(data)
            .map(([key, value]) => `${key}:${typeof value}/${value?.constructor?.name || 'none'}`)
            .join(',')
        : ''
    console.error(
      `[ws] type=${typeof data} constructor=${data?.constructor?.name || 'none'} keys=${keys} values=${valueTypes} length=${data?.length ?? data?.byteLength ?? 'n/a'}`,
    )
  }

  if (data && typeof data === 'object' && 'payload' in data) {
    data = data.payload
  }

  if (Buffer.isBuffer(data)) {
    data = data.toString('utf8')
  } else if (data instanceof Uint8Array) {
    data = Buffer.from(data).toString('utf8')
  } else if (data instanceof ArrayBuffer) {
    data = Buffer.from(data).toString('utf8')
  }

  if (typeof data !== 'string') {
    malformedFrames += 1
    return
  }

  try {
    const message = JSON.parse(data)
    const keys = message && typeof message === 'object' ? Object.keys(message) : []
    if (
      !message ||
      typeof message !== 'object' ||
      typeof message.event !== 'string' ||
      keys.some((key) => key !== 'event' && key !== 'payload')
    ) {
      malformedFrames += 1
      return
    }
    events.push(message.event)
  } catch {
    malformedFrames += 1
  }
}

async function writeOutput(result) {
  const serialized = JSON.stringify(result, null, 2)
  process.stdout.write(`${serialized}\n`)
  if (outputPath) {
    await writeFile(outputPath, `${serialized}\n`, 'utf8')
  }
}

async function captureFailure() {
  await mkdir(artifactDir, { recursive: true })
  if (page) {
    await page.screenshot({ path: `${artifactDir}/failure.png`, fullPage: true }).catch(() => {})
  }
}

try {
  const target = new URL(baseURL)
  browser = await chromium.launch({
    headless: process.env.NEKO_E2E_HEADLESS !== '0',
    args: ['--no-sandbox', '--use-fake-ui-for-media-stream', '--use-fake-device-for-media-stream', '--autoplay-policy=no-user-gesture-required'],
  })
  const context = await browser.newContext({
    viewport: { width: viewportWidth, height: viewportHeight },
    permissions: ['microphone'],
  })
  page = await context.newPage()
  page.setDefaultTimeout(timeout)
  page.on('websocket', (socket) => {
    socket.on('framereceived', parseFrame)
    socket.on('framesent', parseFrame)
  })

  const loginResponse = page.waitForResponse(
    (response) => response.url().endsWith('/api/login') && response.request().method() === 'POST',
  )
  await page.goto(target.toString(), { waitUntil: 'domcontentloaded' })
  await page.getByTestId('displayname-input').fill(username)
  await page.getByTestId('password-input').fill(password)
  await page.getByTestId('connect-submit').click()

  const login = await loginResponse
  if (!login.ok()) {
    throw new Error(`login failed with HTTP ${login.status()}`)
  }

  await page.getByTestId('connection-indicator').waitFor({ state: 'attached' })
  await page.waitForFunction(
    () => document.querySelector('[data-testid="connection-indicator"]')?.classList.contains('connected'),
    undefined,
    { timeout },
  )
  const connectionMs = Math.round(performance.now() - startedAt)

  await page.waitForFunction(
    () => {
      const video = document.querySelector('[data-testid="remote-video"]')
      return video && video.readyState >= 2 && video.videoWidth > 0 && video.videoHeight > 0
    },
    undefined,
    { timeout },
  )

  const frame = await page.evaluate(async () => {
    const video = document.querySelector('[data-testid="remote-video"]')
    if (!(video instanceof HTMLVideoElement)) {
      throw new Error('remote video element is missing')
    }
    await video.play().catch(() => {})
    if (typeof video.requestVideoFrameCallback === 'function') {
      await new Promise((resolve) => video.requestVideoFrameCallback(() => resolve()))
    } else {
      await new Promise((resolve) => requestAnimationFrame(() => resolve()))
    }
    return { width: video.videoWidth, height: video.videoHeight, readyState: video.readyState }
  })
  const firstFrameMs = Math.round(performance.now() - startedAt)

  if (process.env.NEKO_E2E_VOICE === '1') {
    const callers = []
    for (let i = 0; i < 6; i++) {
      const callContext = await browser.newContext({ permissions: ['microphone'] })
      await callContext.addInitScript(() => {
        window.voiceTestPeers = []
        const Original = window.RTCPeerConnection
        window.RTCPeerConnection = class extends Original {
          constructor(...args) { super(...args); window.voiceTestPeers.push(this) }
        }
      })
      const caller = await callContext.newPage()
      caller.setDefaultTimeout(timeout)
      await caller.goto(target.toString())
      await caller.getByTestId('displayname-input').fill(`voice-test-${i}`)
      await caller.getByTestId('password-input').fill(password)
      await caller.getByTestId('connect-submit').click()
      await caller.locator('button[aria-controls="room-panel"]').click()
      await caller.getByTestId('voice-join').click()
      await caller.getByTestId('voice-leave').waitFor()
      callers.push(caller)
    }
    for (const caller of callers) {
      await caller.waitForFunction(async () => {
        let receiving = 0
        for (const peer of window.voiceTestPeers) {
          // Desktop peer has video receivers; count independent audio peers only.
          if (peer.getReceivers().some(receiver => receiver.track.kind === 'video')) continue
          const stats = await peer.getStats()
          if ([...stats.values()].some(stat => stat.type === 'inbound-rtp' && stat.kind === 'audio' && stat.packetsReceived > 0)) receiving++
        }
        return receiving === 5
      }, undefined, { timeout })
    }
    await callers[0].getByTestId('voice-mute').click()
    await callers[0].waitForFunction(() => window.voiceTestPeers.filter(peer => !peer.getReceivers().some(receiver => receiver.track.kind === 'video')).every(peer => peer.getSenders().every(sender => !sender.track || !sender.track.enabled)))
    await callers[0].getByTestId('voice-leave').click()
    await callers[1].waitForFunction(() => document.querySelectorAll('.voice-call li').length === 5)
    for (const caller of callers) await caller.context().close()
  }

  if (process.env.NEKO_E2E_UI === '1') {
    const errors = []
    page.on('pageerror', (error) => errors.push(error.message))
    const toolbar = page.locator('.video-menu')
    await toolbar.locator('button').filter({ has: page.locator('.fa-desktop') }).click()
    await page.locator('.video .context li').first().waitFor({ state: 'visible' })
    await page.locator('[data-testid="remote-video"]').click({ force: true })
    await toolbar.locator('button').filter({ has: page.locator('.fa-expand') }).click()
    await page.waitForFunction(() => !!document.fullscreenElement)
    await page.evaluate(() => document.exitFullscreen())
    await toolbar.locator('button').filter({ has: page.locator('.fa-expand') }).waitFor({ state: 'visible' })
    await page.setViewportSize({ width: 390, height: 844 })
    await page.getByTestId('remote-video').waitFor({ state: 'visible' })
    if (errors.length) throw new Error(`UI errors: ${errors.join('; ')}`)
  }

  if (!events.includes('system/init')) {
    throw new Error('signaling did not deliver system/init')
  }
  if (malformedFrames > 0) {
    throw new Error(`signaling delivered ${malformedFrames} malformed envelope(s)`)
  }
  if (events.includes('screen_sizes_list')) {
    throw new Error('deprecated screen_sizes_list event was observed')
  }

  const result = {
    baseURL: target.origin,
    username,
    connectionMs,
    firstFrameMs,
    video: frame,
    viewport: { width: viewportWidth, height: viewportHeight },
    signalingEvents: [...new Set(events)],
    malformedFrames,
    browser: (await browser.version()).trim(),
    profile: process.env.NEKO_E2E_PROFILE || 'unspecified',
  }
  await writeOutput(result)
} catch (error) {
  await captureFailure()
  const message = error instanceof Error ? error.message : String(error)
  await writeOutput({
    baseURL,
    username,
    error: message,
    signalingEvents: [...new Set(events)],
    malformedFrames,
  })
  process.exitCode = 1
} finally {
  await browser?.close()
}
