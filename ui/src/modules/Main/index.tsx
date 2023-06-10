import React, { useCallback, useEffect, useState } from "react";
import { Col, Container, Nav, Navbar, Row } from "react-bootstrap";
import './style.css'
import { Outlet, useLocation } from "react-router-dom";

const Main = () => {
  const location = useLocation()
  

  return (
    <>
      <Navbar bg="dark" sticky="top">
        <Container fluid>
          <Navbar.Brand href="#home" style={{color: "#fff"}}>Cloudx</Navbar.Brand>
        </Container>
      </Navbar>
      <Container fluid className="main">
        <Row>
          <Col 
            md={3}
            lg={2}
            className="d-md-block bg-light sidebar"
          >
            <Nav activeKey={location.pathname} className="flex-column pt-3">
              <Nav.Link href="/">Shared infras</Nav.Link>
              <Nav.Link href="/connection-interfaces">Connection Interfaces</Nav.Link>
              <Nav.Link href="/cloud-accounts">Cloud accounts</Nav.Link>
              <Nav.Link href="/plugins">Plugins</Nav.Link>
            </Nav>
          </Col>
          <Col md={9} lg={10} className="px-md-4 ms-sm-auto">
            <Outlet />
          </Col>
        </Row>
      </Container>
    </>
  )
}

export default Main