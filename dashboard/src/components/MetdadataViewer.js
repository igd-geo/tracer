import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';


const useStyles = makeStyles(() => ({
  paper: {
    textAlign: 'center',
    height: '75vh',
  },
}));

export default function MetadataViewer() {
  const classes = useStyles();

  return (
    <Grid
      container
      justify="center"
      alignItems="stretch"
      spacing={3}
    >
      <Grid item xs={6}>
        <Paper className={classes.paper}>Old</Paper>
      </Grid>
      <Grid item xs={6}>
        <Paper className={classes.paper}>New</Paper>
      </Grid>
    </Grid>
  );
}
