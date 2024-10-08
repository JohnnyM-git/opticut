// import React from "react";
import { FunctionComponent, useState } from "react";
import {
  Box,
  Button,
  Divider,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  TextField,
} from "@mui/material";
import { FiberNew, MonitorHeart, Settings, Storage } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import "./styles/navbar.css";
import { useNavigate } from "react-router-dom";

interface NavbarProps {
  // startbackend?: any;
  checkForUpdates?: any;
  // stopbackend?: any;
}

export const Navbar: FunctionComponent<NavbarProps> = ({
  // startbackend,
  checkForUpdates,
}) => {
  const [open, setOpen] = useState<boolean>(false);
  const [jobId, setJobId] = useState("");
  const navigate = useNavigate();
  const DrawerList = (
    <Box
      sx={{ width: 250 }}
      role="presentation"
      onClick={() => toggleDrawer(false)}
    >
      <List>
        <ListItem disablePadding onClick={() => navigate("/")}>
          <ListItemButton>
            <ListItemIcon>
              <FiberNew />
            </ListItemIcon>
            <ListItemText primary={"New Job"} />
          </ListItemButton>
        </ListItem>
        <Divider />

        <ListItem disablePadding onClick={() => navigate("/localjobs")}>
          <ListItemButton>
            <ListItemIcon>
              <Storage />
            </ListItemIcon>
            <ListItemText primary={"Local Jobs"} />
          </ListItemButton>
        </ListItem>
        <Divider />

        {/*<ListItem disablePadding>*/}
        {/*  <ListItemButton>*/}
        {/*    <ListItemIcon>*/}
        {/*      <Cloud />*/}
        {/*    </ListItemIcon>*/}
        {/*    <ListItemText primary={"Cloud Jobs"} />*/}
        {/*  </ListItemButton>*/}
        {/*</ListItem>*/}
        {/*<Divider />*/}

        {/*<ListItem disablePadding>*/}
        {/*  <ListItemButton>*/}
        {/*    <ListItemIcon>*/}
        {/*      <People />*/}
        {/*    </ListItemIcon>*/}
        {/*    <ListItemText primary={"Customers"} />*/}
        {/*  </ListItemButton>*/}
        {/*</ListItem>*/}
        {/*<Divider />*/}

        <ListItem disablePadding onClick={() => navigate("/settings")}>
          <ListItemButton>
            <ListItemIcon>
              <Settings />
            </ListItemIcon>
            <ListItemText primary={"Settings"} />
          </ListItemButton>
        </ListItem>
        <Divider />

        <ListItem disablePadding onClick={() => navigate("/healthstatus")}>
          <ListItemButton>
            <ListItemIcon>
              <MonitorHeart />
            </ListItemIcon>
            <ListItemText primary={"Health Status"} />
          </ListItemButton>
        </ListItem>
        <Divider />
      </List>
    </Box>
  );

  const toggleDrawer = (openState: boolean): void => {
    setOpen(openState);
  };

  return (
    <div className={"navbar"}>
      <div className={"menu"}>
        <Button onClick={() => toggleDrawer(true)}>
          <MenuIcon />
        </Button>
        <Drawer
          className={"drawer"}
          open={open}
          onClose={() => toggleDrawer(false)}
        >
          {DrawerList}
        </Drawer>
      </div>
      <div className={"search"}>
        <TextField
          value={jobId}
          placeholder={"Enter Job Number"}
          onChange={(e) => setJobId(e.target.value)}
        />
        {/*<Button onClick={() => startbackend()}>Start</Button>*/}
        <Button onClick={() => checkForUpdates()}>Check</Button>
        <Button onClick={() => navigate(`/results/${jobId}`)}>Go</Button>
      </div>
    </div>
  );
};
