import {
  SET_ACTIVE_NODE,
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
    { source: 'Travis', target: 'Rake' },
    { source: 'Diana', target: 'Rake' },
    { source: 'Diana', target: 'Rachel' },
    { source: 'Rachel', target: 'Rake' },
    { source: 'Rachel', target: 'Shawn' },
    { source: 'Emerald', target: 'Rachel' },
  ],
  activeNode: { name: '' },
};

const graph = (state = initialState, action) => {
  switch (action.type) {
    case SET_ACTIVE_NODE:
    {
      console.log(initialState);
      const node = initialState.nodes.find((n) => n.id === action.payload);
      return { ...state, activeNode: node };
    }
    default: {
      return state;
    }
  }
};

export default graph;
