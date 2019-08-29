import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

function DetailViewer(props) {
  const { activeNode } = props;
  const details = Object.keys(activeNode).map((key) => (
    <>
      <li>{key}</li>
      <li>{activeNode[key]}</li>
    </>
  ));
  return (
    <div>
      {details}
    </div>
  );
}
const mapStateToProps = (state) => ({
  activeNode: state.graph.activeNode,
});

DetailViewer.propTypes = {
  activeNode: PropTypes.any,
};

export default connect(mapStateToProps, null)(DetailViewer);
