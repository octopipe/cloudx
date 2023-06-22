import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Badge, Button, Card, Col, Container, ListGroup, Row, Tab, Tabs } from "react-bootstrap";
import { useParams } from "react-router-dom";
import SharedInfraViewDiagram from "./Diagram";
import { toEdges, toNodes } from "./utils";
import "./index.css"

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
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3">
          <h1 className="h2">{sharedInfra?.name}</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
        <Tabs
          defaultActiveKey="overview"
          id="fill-tab-example"
          className="mb-3"
        >
          <Tab eventKey="overview" title="Overview">
            <SharedInfraViewDiagram
              initialNodes={sharedInfra?.plugins ? toNodes(sharedInfra.plugins) : []}
              initialEdges={sharedInfra?.plugins ? toEdges(sharedInfra.plugins) : []}
            />
          </Tab>
          <Tab eventKey="executions" title="Executions">
            <div>
              <Accordion>
                {sharedInfra?.status?.executions?.map((item: any, idx: any) => (
                  <Accordion.Item eventKey={idx}>
                    <Accordion.Header>Execution #{ idx }</Accordion.Header>
                    <Accordion.Body>
                      <SharedInfraViewDiagram
                        initialNodes={item?.plugins ? toNodes(item.plugins) : []}
                        initialEdges={item?.plugins ? toEdges(item.plugins) : []}
                      />
                    </Accordion.Body>
                  </Accordion.Item>
                ))}
              </Accordion>
            </div>
          </Tab>
        </Tabs>
      </Container>
      <Container fluid>
        
        {/* <div>
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
        </div> */}

        
      </Container>
    </>
  )
}

export default SharedInfraView