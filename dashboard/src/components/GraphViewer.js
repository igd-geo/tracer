import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import * as d3 from 'd3';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { setActiveNode } from '../redux/actions';

const useStyles = makeStyles(() => ({
  paper: {
    textAlign: 'center',
    height: '75vh',
  },
}));

function GraphViewer(props) {
  const radius = 30;
  const classes = useStyles();
  const { nodes, edges, dispatchActiveNode } = props;


  useEffect(() => {
    const svg = d3.select('.graphSVG');

    if (!svg.selectAll('*').empty()) {
      svg.selectAll('*').remove();
    }
    const { height, width } = svg.node().getBoundingClientRect();
    const graph = {
      nodes,
      edges,
    };

    const simulation = d3.forceSimulation()
      .nodes(graph.nodes);
    simulation
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collide', d3.forceCollide().radius(50));
    simulation.force('links', d3.forceLink(edges).id((d) => d.name));

    const link = svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(graph.edges)
      .enter()
      .append('line');

    const node = svg.append('g')
      .attr('class', 'nodes')
      .selectAll('circle')
      .data(graph.nodes)
      .enter()
      .append('circle')
      .attr('r', radius)
      .attr('fill', 'red')
      .on('click', (clickedNode) => {
        dispatchActiveNode(clickedNode.id);
      })
      .on('mouseover', () => {
        d3.select(d3.event.target).transition().duration(50).attr('fill', 'orange');
        d3.select(d3.event.target).style('cursor', 'pointer');
      })
      .on('mouseout', () => {
        d3.select(d3.event.target).style('cursor', 'default');
        d3.select(d3.event.target).transition().duration(50).attr('fill', 'red');
      });

    const tickActions = () => {
      node
        .attr('cx', (d) => d.x)
        .attr('cy', (d) => d.y);
      link
        .attr('x1', (d) => d.source.x)
        .attr('y1', (d) => d.source.y)
        .attr('x2', (d) => d.target.x)
        .attr('y2', (d) => d.target.y);
    };
    simulation.on('tick', tickActions);

    return () => {
      svg.selectAll('*').remove();
    };
  }, [nodes, edges, dispatchActiveNode]);

  return (
    <Grid
      container
      justify="center"
      alignItems="stretch"
      spacing={3}
    >
      <Grid item xs={12}>
        <Paper className={classes.paper}>
          <svg className="graphSVG" height="800px" width="100%" />
        </Paper>
      </Grid>
    </Grid>
  );
}

const mapStateToProps = (state) => ({
  nodes: state.graph.nodes,
  edges: state.graph.edges,
  activeNode: state.graph.activeNode,
});

const mapDispatchToProps = (dispatch) => ({
  dispatchActiveNode: (node) => dispatch(setActiveNode(node)),
});


GraphViewer.propTypes = {
  nodes: PropTypes.any,
  edges: PropTypes.any,
  dispatchActiveNode: PropTypes.func,
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(GraphViewer);
