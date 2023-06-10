import React from "react";
import { Button, Container } from "react-bootstrap";

const CloudAccounts = () => {
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Connections Interfaces</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
      </Container>
    </>
  )
}

export default CloudAccounts