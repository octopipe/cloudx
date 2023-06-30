import React, { useCallback, useEffect, useState } from "react";
import { Col, Container, Nav, Navbar, Row } from "react-bootstrap";
import './style.css'
import { NavLink, Outlet, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const Main = () => {
  const location = useLocation()
  

  return (
    <>
      <Navbar bg="dark" sticky="top">
        <Container fluid>
          <Navbar.Brand href="#home" style={{color: "#fff"}}>Cloudx</Navbar.Brand>
        </Container>
      </Navbar>
      <div>
        <div className="d-md-block main__sidebar">
          <Nav activeKey={location.pathname} className="flex-column pt-3">
            <NavLink className="nav-link" to="/shared-infras"><FontAwesomeIcon icon="layer-group" /></NavLink>
            <NavLink className="nav-link" to="/connection-interfaces"><FontAwesomeIcon icon="diagram-project" /></NavLink>
            <NavLink className="nav-link" to="/providers-config"><FontAwesomeIcon icon="cloud" /></NavLink>
          </Nav>
        </div>
        <div className="main__content">
          <Outlet />
        </div>
      </div>
    </>
  )
}

export default Main