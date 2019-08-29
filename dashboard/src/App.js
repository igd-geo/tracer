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
import TextField from '@material-ui/core/TextField';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import { makeStyles } from '@material-ui/core/styles';
import MetadataViewer from './components/MetdadataViewer';
import GraphViewer from './components/GraphViewer';
import DetailViewer from './components/DetailViewer';

const drawerWidth = 420;
const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  toolbar: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    height: '100vh',
    oveflow: 'auto',
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4),
  },
  textField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 200,
  },
}));


function TabPanel(props) {
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
      <Box p={3}>{children}</Box>
    </Typography>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

function App() {
  const classes = useStyles();
  const [value, setValue] = React.useState(0);
  const [values, setValues] = React.useState({
    name: '',
  });

  function handleTabChange(event, newValue) {
    setValue(newValue);
  }

  const handleIdChange = (name) => (event) => {
    setValues({ ...values, [name]: event.target.value });
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
          <Grid container direction="column" spacing={3}>
            <Grid item>
              <Toolbar>
                <TextField
                  id="standard-name"
                  label="ID"
                  className={classes.textField}
                  value={values.name}
                  onChange={handleIdChange('name')}
                  margin="normal"
                  variant="filled"
                />
                <Tabs value={value} onChange={handleTabChange} aria-label="simple tabs example">
                  <Tab label="Graph" {...a11yProps(0)} />
                  <Tab label="Metadata" {...a11yProps(1)} />
                </Tabs>
              </Toolbar>
            </Grid>
            <Grid item>
              <TabPanel value={value} index={0}>
                <GraphViewer />
              </TabPanel>
              <TabPanel value={value} index={1}>
                <MetadataViewer />
              </TabPanel>
            </Grid>
          </Grid>
        </Container>
      </main>
    </div>
  );
}

export default App;
