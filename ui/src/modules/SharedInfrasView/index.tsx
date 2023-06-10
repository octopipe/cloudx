import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Card, Col, Container, ListGroup, Row } from "react-bootstrap";
import { useParams } from "react-router-dom";

const SharedInfraView = () => {
  const { name } = useParams()
  const [sharedInfra, setSharedInfra] = useState<any>()

  const getSharedInfra = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const item = await res.json()
    console.log(item)
    setSharedInfra(item)
  }, [])

  useEffect(() => {
    if (!name)
      return

    getSharedInfra(name)
  }, [])
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">{sharedInfra?.metadata?.name}</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
      </Container>
      <Container fluid>
        <div>
          <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-1">
            <h1 className="h5">Plugins</h1>
          </div>
          {sharedInfra?.spec?.plugins?.map((item: any, idx: any) => (
            <Card className="mb-2">
              <Card.Body>
                <Card.Title>{item?.name}</Card.Title>
                <Card.Subtitle className="mb-2 text-muted">{item?.ref}</Card.Subtitle>
              </Card.Body>
            </Card>
          ))}
        </div>

        <div>
          <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-1">
            <h1 className="h5">Executions</h1>
          </div>
          {sharedInfra?.status?.executions?.map((item: any, idx: any) => (
            <Card className="mb-2">
              <Card.Body>
                <Card.Title>{item?.status}</Card.Title>
                <Card.Subtitle className="mb-2 text-muted">{item?.ref}</Card.Subtitle>
              </Card.Body>
            </Card>
          ))}
        </div>
      </Container>
    </>
  )
}

export default SharedInfraView