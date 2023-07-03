import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Card, Col, Container, ListGroup, Row } from "react-bootstrap";
import { Link, useNavigate } from "react-router-dom";

const ProviderConfig = () => {
  const navigate = useNavigate()
  const [list, setList] = useState<any>()

  const getList = useCallback(async () => {
    const res = await fetch("http://localhost:8080/providers-configs")
    const list = await res.json()

    setList(list)
  }, [])

  useEffect(() => {
    getList()
  }, [])
  
  return (
    <div  style={{margin: "80px"}}>
      <div className="d-flex justify-content-between my-4">
        <h1 className="h2">Providers config</h1>
        <Button onClick={() => navigate('/shared-infras/create')}>
          <FontAwesomeIcon icon="add" /> Create
        </Button>  
      </div>
      <Card>
        <Card.Body>
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
    </div>
  )
}

export default ProviderConfig