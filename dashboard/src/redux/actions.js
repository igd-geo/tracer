import {
  SET_ACTIVE_NODE,
} from './actionTypes';

export const setActiveNode = (id) => ({
  type: SET_ACTIVE_NODE,
  payload: id,
});
