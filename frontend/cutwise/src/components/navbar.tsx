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
import { Home, Inbox, Mail } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { getJobData } from "../functions/getJobData.ts";
import "./styles/navbar.css";

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
  const DrawerList = (
    <Box
      sx={{ width: 250 }}
      role="presentation"
      onClick={() => toggleDrawer(false)}
    >
      <List>
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemIcon>
              <Home />
            </ListItemIcon>
            <ListItemText primary={"Home / All Jobs"} />
          </ListItemButton>
        </ListItem>
        <Divider />

        <ListItem disablePadding>
          <ListItemButton>
            <ListItemIcon>
              <Home />
            </ListItemIcon>
            <ListItemText primary={"Customers"} />
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
        <Button onClick={() => getJobData(jobId)}>Go</Button>
      </div>
    </div>
  );
};
