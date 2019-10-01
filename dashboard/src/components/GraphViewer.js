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
    overflow: 'hidden',
    height: '100%',
  },
  grid: {
    height: '100%',
  },
  gridContainer: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
}));

const graphWidth = '100%';
const graphHeight = '100%';
const radius = 30;

function GraphViewer(props) {
  const { nodes, edges, dispatchActiveNode } = props;
  const classes = useStyles();

  const defaultColor = 'black';
  const highlightColor = 'orange';

  const rootColor = 'blue';
  const entityColor = 'yellow';
  const activityColor = 'green';
  const agentColor = 'red';

  const nodeColor = (type) => {
    switch (type) {
      case 'root':
        return rootColor;
      case 'entity':
        return entityColor;
      case 'activity':
        return activityColor;
      case 'agent':
        return agentColor;
      default:
        return defaultColor;
    }
  };

  const edgeColor = (type) => {
    switch (type) {
      case 'wasDerivedFrom':
        return entityColor;
      case 'wasGeneratedBy':
        return activityColor;
      case 'wasAssociatedWith':
        return agentColor;
      case 'wasAttributedTo':
        return agentColor;
      case 'actedOnBehalfOf':
        return agentColor;
      case 'used':
        return entityColor;
      default:
        return defaultColor;
    }
  };

  useEffect(() => {
    const svg = d3.select('.graphSVG');

    if (!svg.selectAll('*').empty()) {
      svg.selectAll('*').remove();
    }
    const { height, width } = svg.node().getBoundingClientRect();


    const simulation = d3.forceSimulation()
      .nodes(nodes)
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collide', d3.forceCollide().radius(4 * radius))
      .force('links', d3.forceLink(edges).id((d) => d.id));

    const link = svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(edges)
      .enter()
      .append('line')
      .style('stroke', (d) => edgeColor(d.edgeType));

    const node = svg.append('g')
      .attr('class', 'nodes')
      .selectAll('circle')
      .data(nodes)
      .enter()
      .append('circle')
      .attr('r', radius)
      .attr('fill', (d) => nodeColor(d.nodeType))
      .on('click', (clickedNode) => {
        dispatchActiveNode(clickedNode.id);
      })
      .on('mouseover', () => {
        d3.select(d3.event.target).transition().duration(50).attr('fill', highlightColor);
        d3.select(d3.event.target).style('cursor', 'pointer');
      })
      .on('mouseout', () => {
        d3.select(d3.event.target).style('cursor', 'default');
        d3.select(d3.event.target).transition().duration(50).attr('fill', (d) => nodeColor(d.nodeType));
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
      className={classes.gridContainer}
    >
      <Grid item xs={12} className={classes.grid}>
        <Paper className={classes.paper}>
          <svg className="graphSVG" height={graphHeight} width={graphWidth} />
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
  nodes: PropTypes.array,
  edges: PropTypes.array,
  dispatchActiveNode: PropTypes.func,
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(GraphViewer);
