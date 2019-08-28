import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import * as d3 from 'd3';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
  paper: {
    textAlign: 'center',
    height: '75vh',
  },
}));

const nodes_data = [
  { name: 'Travis', sex: 'M' },
  { name: 'Rake', sex: 'M' },
  { name: 'Diana', sex: 'F' },
  { name: 'Rachel', sex: 'F' },
  { name: 'Shawn', sex: 'M' },
  { name: 'Emerald', sex: 'F' },
];
const links_data = [
  { source: 'Travis', target: 'Rake' },
  { source: 'Diana', target: 'Rake' },
  { source: 'Diana', target: 'Rachel' },
  { source: 'Rachel', target: 'Rake' },
  { source: 'Rachel', target: 'Shawn' },
  { source: 'Emerald', target: 'Rachel' },
];


function GraphViewer(props) {
  const classes = useStyles();

  useEffect(() => {
    const svg = d3.select('.graphSVG');
    if (!svg.selectAll('*').empty()) {
      return () => {
        svg.selectAll('*').remove();
      };
    }
    const { height, width } = svg.node().getBoundingClientRect();
    const link_force = d3.forceLink(links_data)
      .id((d) => d.name);

    const simulation = d3.forceSimulation()
      .nodes(nodes_data);
    simulation
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collide', d3.forceCollide().strength(0.2).radius(30).iterations(1));
    simulation.force('links', link_force);

    const link = svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(links_data)
      .enter()
      .append('line');

    const node = svg.append('g')
      .attr('class', 'nodes')
      .selectAll('circle')
      .data(nodes_data)
      .enter()
      .append('circle')
      .attr('r', 20)
      .attr('fill', 'red');

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
  });

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

GraphViewer.propTypes = {
  classes: PropTypes.any,
};

export default GraphViewer;
