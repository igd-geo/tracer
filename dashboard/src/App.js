import React from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Drawer from '@material-ui/core/Drawer';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Box from '@material-ui/core/Box';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import { makeStyles } from '@material-ui/core/styles';
import InputBase from '@material-ui/core/InputBase';
import SearchIcon from '@material-ui/icons/Search';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import { connect } from 'react-redux';
import MetadataViewer from './components/MetdadataViewer';
import GraphViewer from './components/GraphViewer';
import DetailViewer from './components/DetailViewer';
import { setGraph } from './redux/actions';


const drawerWidth = '33vh';
const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    height: '100vh',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
    overflow: 'hidden',
  },
  drawerPaper: {
    width: drawerWidth,
    position: 'static',
    overflow: 'hidden',
  },
  toolbar: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    oveflow: 'auto',
    overflow: 'hidden',
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4),
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    marginLeft: theme.spacing(1),
    flex: 1,
  },
  iconButton: {
    padding: 10,
  },
  searchRoot: {
    padding: '2px 4px',
    display: 'flex',
    alignItems: 'center',
    width: 200,
    marginRight: theme.spacing(5),
  },
  flexContent: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
  },
  tab: {
    flexGrow: 1,
  },
  tabContent: {
    height: '100%',
  },
  tabBox: {
    height: '100%',
    paddingBottom: '50px',
  },
}));

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

function TabPanel(props) {
  const classes = useStyles();
  const {
    children, value, index, ...other
  } = props;

  return (
    <Typography
      component="div"
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      <Box p={3} className={classes.tabBox}>
        {children}
        {' '}
      </Box>
    </Typography>
  );
}

function App(props) {
  const classes = useStyles();
  const [value, setValue] = React.useState(0);
  const [id, setID] = React.useState('');
  const { dispatchGraph } = props;

  const handleTabChange = (event, newValue) => setValue(newValue);

  const handleInputChange = (event) => {
    setID(event.target.value);
  };

  const buildGraph = (data) => ({
    nodes: data.nodes,
    edges: data.edges,
  });

  const handleClick = () => {
    fetch('http://localhost:1234/api?query={graph(id:"entity_123456789"){json nodes edges}}', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
    })
      .then((r) => r.json())
      .then((body) => {
        const graph = buildGraph(body.data.graph);
        dispatchGraph(graph);
      });

    // used for search field
    console.log(id);
  };

  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppBar position="absolute" className={classes.appBar}>
        <Toolbar>
          <Typography variant="h5" noWrap>
            Tracer Dashboard
          </Typography>
        </Toolbar>
      </AppBar>
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{
          paper: classes.drawerPaper,
        }}
        anchor="left"
      >
        <div className={classes.toolbar} />
        <DetailViewer />
      </Drawer>
      <main className={classes.content}>
        <div className={classes.toolbar} />
        <Container maxWidth="xl" className={classes.container}>
          <Grid container direction="column" spacing={3} className={classes.flexContent}>
            <Grid item>
              <Toolbar>
                <Paper className={classes.searchRoot}>
                  <InputBase
                    className={classes.input}
                    placeholder="Document ID"
                    inputProps={{ 'aria-label': 'search document id' }}
                    onChange={handleInputChange}
                  />
                  <IconButton className={classes.iconButton} aria-label="search" onClick={handleClick}>
                    <SearchIcon />
                  </IconButton>
                </Paper>
                <Tabs value={value} onChange={handleTabChange} aria-label="simple tabs example">
                  <Tab label="Graph" {...a11yProps(0)} />
                  <Tab label="Metadata" {...a11yProps(1)} />
                </Tabs>
              </Toolbar>
            </Grid>
            <Grid item className={classes.flexContent}>
              <TabPanel value={value} index={0} className={classes.tab}>
                <GraphViewer className={classes.tabContent} />
              </TabPanel>
              <TabPanel value={value} index={1} className={classes.tab}>
                <MetadataViewer className={classes.tabContent} />
              </TabPanel>
            </Grid>
          </Grid>
        </Container>
      </main>
    </div>
  );
}


const mapDispatchToProps = (dispatch) => ({
  dispatchGraph: (graph) => dispatch(setGraph(graph)),
});

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

App.propTypes = {
  dispatchGraph: PropTypes.func,
};

export default connect(
  null,
  mapDispatchToProps,
)(App);
