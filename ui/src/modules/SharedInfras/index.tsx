import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Card, Col, Container, ListGroup, Row } from "react-bootstrap";
import { Link } from "react-router-dom";

const SharedInfras = () => {
  const [list, setList] = useState<any>()

  const getList = useCallback(async () => {
    const res = await fetch("http://localhost:8080/shared-infras")
    const list = await res.json()

    setList(list)
  }, [])

  useEffect(() => {
    getList()
  }, [])
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Shared infras</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
      </Container>
      <ListGroup as="ul">
        {list?.items?.map((item: any, idx: any) => (
        <ListGroup.Item
          as="li"
          className="d-flex justify-content-between align-items-start"
        >
          <div className="ms-2 me-auto">
            <div className="fw-bold">
              <Link to={item?.name}>{item?.name}</Link>
            </div>
            {item?.description}
          </div>
          <Badge bg="primary" pill>
            {item?.spec?.plugins?.length}
          </Badge>
        </ListGroup.Item>
        ))}
      </ListGroup>
    </>
  )
}

export default SharedInfras