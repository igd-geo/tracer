import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

function DetailViewer(props) {
  const { activeNode } = props;
  const blackList = ['x', 'y', 'vx', 'vy', 'index', 'nodeType', 'edgeType'];
  return Object.keys(activeNode).map((key) => {
    if (blackList.includes(key)) { return null; }
    return (
      <Fragment key={key}>
        <li>
          {key}
          :
          {' '}
          {activeNode[key]}
        </li>
      </Fragment>
    );
  });
}
const mapStateToProps = (state) => ({
  activeNode: state.graph.activeNode,
});

DetailViewer.propTypes = {
  activeNode: PropTypes.any,
};

export default connect(mapStateToProps, null)(DetailViewer);
