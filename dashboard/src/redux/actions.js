import {
  SET_ACTIVE_NODE,
  SET_GRAPH,
} from './actionTypes';

export const setActiveNode = (id) => ({
  type: SET_ACTIVE_NODE,
  payload: id,
});

export const setGraph = (graph) => ({
  type: SET_GRAPH,
  payload: graph,
});
