/**
 * Small, framework-neutral definitions used by the existing store modules.
 *
 * The former store helpers only described module definitions. Keeping the
 * module definitions declarative lets the Pinia adapter in index.ts expose
 * the same public store API while the components are migrated independently.
 */
export const getterTree = <S, G>(_state: () => S, getters: G): G => getters

export const mutationTree = <S, M>(_state: () => S, mutations: M): M => mutations

export const actionTree = <T, A>(_tree: T, actions: A): A => actions
