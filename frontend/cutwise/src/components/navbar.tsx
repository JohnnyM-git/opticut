// import React from "react";
import React, { FunctionComponent, useEffect, useState } from "react";
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
import {
  Cloud,
  DataArray,
  DataObject,
  Dataset,
  DataUsage,
  Home,
  Inbox,
  LocalActivity,
  LocalAtm,
  LocalBar,
  Mail,
  People,
  Settings,
  Storage,
} from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { getJobData } from "../functions/getJobData.ts";
import "./styles/navbar.css";
import { useNavigate } from "react-router-dom";

interface NavbarProps {
  jobId: string;
  setJobData?: any;
  setJobId?: any;
}

export const Navbar: FunctionComponent<NavbarProps> = ({
  jobId,
  setJobData,
  setJobId,
}) => {
  const [open, setOpen] = useState<boolean>(false);
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
              <Home />
            </ListItemIcon>
            <ListItemText primary={"Home"} />
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
        <Drawer open={open} onClose={() => toggleDrawer(false)}>
          {DrawerList}
        </Drawer>
      </div>
      <div className={"search"}>
        <TextField
          value={jobId}
          placeholder={"Enter Job Number"}
          onChange={(e) => setJobId(e.target.value)}
        />
        <Button onClick={() => navigate(`/results/${jobId}`)}>Go</Button>
      </div>
    </div>
  );
};
