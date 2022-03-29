import { GetterTree, ActionTree, MutationTree } from 'vuex'

export const state = () => ({
  currentPageName: '',
})

export type RootState = ReturnType<typeof state>

export const mutations: MutationTree<RootState> = {
  setCurrentPageName: (state, page: string) => (state.currentPageName = page),
}

export const actions: ActionTree<RootState, RootState> = {
  updateCurrentPageName: ({ commit }, page: string) =>
    commit('setCurrentPageName', page),
}

export const getters: GetterTree<RootState, RootState> = {
  getCurrentPageName: (state) => state.currentPageName,
}
