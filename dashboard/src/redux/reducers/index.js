import { combineReducers } from 'redux';
import app from './app';
import metadata from './metadata';
import graph from './graph';

export default combineReducers({
  app, graph, metadata,
});
