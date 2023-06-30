import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Card, Col, Container, ListGroup, Row } from "react-bootstrap";
import { Link, useNavigate } from "react-router-dom";

const SharedInfras = () => {
  const navigate = useNavigate()
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
    <Card style={{margin: "40px"}}>
      <Card.Body>
        <Button
          onClick={() => navigate('/shared-infras/create')}
          style={{ position: 'absolute', right: '20px', borderRadius: '50%'}}>
          <FontAwesomeIcon icon="add" />
        </Button>
        <Card.Title>Shared infras</Card.Title>
        <ListGroup as="ul">
          {list?.items?.map((item: any, idx: any) => (
          <ListGroup.Item
            as="li"
            className="d-flex justify-content-between align-items-start"
          >
            <div className="">
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
      </Card.Body>
    </Card>
  )
}

export default SharedInfras