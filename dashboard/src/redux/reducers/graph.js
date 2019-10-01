import {
  SET_ACTIVE_NODE,
  SET_GRAPH,
} from '../actionTypes';

const initialState = {
  nodes: [
    { id: 'tr', name: 'Travis', sex: 'M' },
    { id: 'ra', name: 'Rake', sex: 'M' },
    { id: 'di', name: 'Diana', sex: 'F' },
    { id: 'rac', name: 'Rachel', sex: 'F' },
    { id: 'sh', name: 'Shawn', sex: 'M' },
    { id: 'em', name: 'Emerald', sex: 'F' },
  ],
  edges: [
    { source: 'tr', target: 'ra' },
    { source: 'di', target: 'ra' },
    { source: 'di', target: 'rac' },
    { source: 'rac', target: 'ra' },
    { source: 'rac', target: 'sh' },
    { source: 'em', target: 'rac' },
  ],
  activeNode: { name: '' },
};

const graph = (state = initialState, action) => {
  switch (action.type) {
    case SET_ACTIVE_NODE:
    {
      const node = state.nodes.find((n) => n.id === action.payload);
      return { ...state, activeNode: node };
    }
    case SET_GRAPH:
    {
      return { ...state, nodes: action.payload.nodes, edges: action.payload.edges };
    }
    default: {
      return state;
    }
  }
};

export default graph;
