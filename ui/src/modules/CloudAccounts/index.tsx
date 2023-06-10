import React from "react";
import { Button, Container } from "react-bootstrap";

const ConnectionInterfaces = () => {
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Cloud Accounts</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
      </Container>
    </>
  )
}

export default ConnectionInterfaces