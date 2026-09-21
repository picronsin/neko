import type { Plugin } from 'vue'
import anime, { StaggerOptions, AnimeTimelineInstance, AnimeParams, AnimeInstance } from 'animejs'

type FunctionBasedParameter = (element: HTMLElement, index: number, length: number) => number
type AnimeTarget = string | object | HTMLElement | SVGElement | NodeList | null
type AnimeFunc = (params: AnimeParams) => AnimeInstance

interface Anime {
  version: string
  speed: number
  running: AnimeInstance[]
  easings: { [EasingFunction: string]: (t: number) => any }
  remove(targets: AnimeTarget | ReadonlyArray<AnimeTarget>): void
  get(targets: AnimeTarget, prop: string): string | number
  path(
    path: string | HTMLElement | SVGElement | null,
    percent?: number,
  ): (prop: string) => {
    el: HTMLElement | SVGElement
    property: string
    totalLength: number
  }
  setDashoffset(el: HTMLElement | SVGElement | null): number
  bezier(x1: number, y1: number, x2: number, y2: number): (t: number) => number
  stagger(value: number | string | ReadonlyArray<number | string>, options?: StaggerOptions): FunctionBasedParameter
  set(targets: AnimeTarget, value: { [AnyAnimatedProperty: string]: any }): void
  // Timeline
  timeline(params?: AnimeParams | ReadonlyArray<AnimeInstance>): AnimeTimelineInstance
  random(min: number, max: number): number
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $anime: any
  }
}

const plugin: Plugin = {
  install(app) {
    app.config.globalProperties.$anime = anime
  },
}

export default plugin
