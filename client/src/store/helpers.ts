/**
 * Small, framework-neutral definitions used by the existing store modules.
 *
 * The former store helpers only described module definitions. Keeping the
 * module definitions declarative lets the Pinia adapter in index.ts expose
 * the same public store API while the components are migrated independently.
 */
export const getterTree = <S, G>(_state: () => S, getters: G): G => getters

export const mutationTree = <S, M extends Record<string, (state: S, ...args: any[]) => unknown>>(
  _state: () => S,
  mutations: M,
): M => mutations

type GetterValues<G> = {
  [K in keyof G]: G[K] extends (...args: any[]) => infer R ? R : never
}

export const actionTree = <
  S,
  G,
  M,
  A extends Record<
    string,
    (
      context: {
        state: S
        getters: GetterValues<G>
      },
      ...args: any[]
    ) => unknown
  >,
>(
  _tree: { state: () => S; getters: G; mutations: M },
  actions: A,
): A => actions
