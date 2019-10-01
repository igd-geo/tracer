import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';

const useStyles = makeStyles(() => ({
  paper: {
    textAlign: 'center',
    overflow: 'scroll',
    height: '100%',
  },
  grid: {
    height: '100%',
  },
  gridContainer: {
    height: '100%',
  },
}));

export default function MetadataViewer() {
  const classes = useStyles();
  const [xml, setXML] = React.useState('');

  useEffect(() => {
    fetch('http://localhost:1234/metadata')
      .then((response) => response.text())
      .then((str) => (setXML(str)));
  });

  return (
    <Grid
      container
      justify="center"
      alignItems="stretch"
      spacing={3}
      className={classes.gridContainer}
    >
      <Grid item xl={6} xs={5} xm={5} className={classes.grid}>
        <Paper className={classes.paper}>{xml}</Paper>
      </Grid>
      <Grid item xl={6} xs={5} xm={5} className={classes.grid}>
        <Paper className={classes.paper}>New</Paper>
      </Grid>
    </Grid>
  );
}
