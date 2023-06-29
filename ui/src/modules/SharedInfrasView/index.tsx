import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useParams } from "react-router-dom";
import SharedInfraViewDiagram from "./Diagram";
import { toEdges, toNodes } from "./utils";
import "./index.css"

const getBadgeVariants = (status: string) => {
  if (status === "RUNNING") {
    return 'primary'
  }

  if (status === "SUCCESS") {
    return 'success'
  }


  return 'danger'
}

const SharedInfraView = () => {
  const { name } = useParams()
  const [sharedInfra, setSharedInfra] = useState<any>()
  const [selectedExecution, setSelectedExecution] = useState<any>()

  const getSharedInfra = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const item = await res.json()

    setSharedInfra(item)
  }, [])


  const getExecution = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/executions/${name}`)
    const item = await res.json()

    setSelectedExecution(item)
  }, [])

  useEffect(() => {
    if (!name)
      return

    const interval = setInterval(() => {
      getSharedInfra(name)
    }, 3000)

    getSharedInfra(name)
    return () => clearInterval(interval)
  }, [])
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2">
          <h1 className="h2">{sharedInfra?.name}</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button variant="secondary">Create</Button>
          </div>
        </div>
        <div className="mb-3">
          {sharedInfra?.description}
        </div>
        <Tabs
          defaultActiveKey="overview"
          id="fill-tab-example"
          className="mb-3"
        >
          <Tab eventKey="overview" title="Overview">
            {/* <h1 className="h4">Plugins</h1> */}
            {/* <Accordion>
              {sharedInfra?.plugins?.map((p: any, idx: any) => (

                <Accordion.Item eventKey={idx}>
                  <Accordion.Header>
                    {p?.name}
                  </Accordion.Header>
                  <Accordion.Body>
                  <div className="mb-1"><strong>Name: </strong>{p?.name}</div>
                  <div className="mb-1"><strong>Ref: </strong>{p?.ref}</div>
                  <div className="mb-1"><strong>Type: </strong>{p?.type}</div>
                  <div>
                    <strong>Inputs: </strong>
                    <Table bordered>
                      <thead>
                        <tr>
                          <th>Name</th>
                          <th>Value</th>
                        </tr>
                      </thead>
                      <tbody>
                        {p?.inputs?.map((i: any) => (
                          <tr>
                            <td>{i.key}</td>
                            <td>{i.value}</td>

                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  </div>
                  </Accordion.Body>

                </Accordion.Item>
              ))}
            </Accordion> */}
            <SharedInfraViewDiagram
              initialNodes={sharedInfra?.plugins ? toNodes(sharedInfra.plugins, "default") : []}
              initialEdges={sharedInfra?.plugins ? toEdges(sharedInfra.plugins) : []}
            />
          </Tab>
          <Tab eventKey="executions" title="Executions">
            <div>
              <Accordion>
                {sharedInfra?.status?.executions?.map((item: any, idx: any) => (
                  <Accordion.Item eventKey={idx}>
                    <Accordion.Header onClick={() => getExecution(item?.name)}>
                      {item?.status === "RUNNING" && (
                        <Spinner animation="border" role="status" variant="primary" size="sm">
                          <span className="visually-hidden">Loading...</span>
                        </Spinner>
                      )}
                      {/* <Badge className="mx-2" bg={getBadgeVariants(selectedExecution?.status?.status)}>{ selectedExecution?.status?.status }</Badge> */}
                      Execution #{ sharedInfra?.status?.executions?.length - idx }
                    </Accordion.Header>
                    <Accordion.Body>
                      {selectedExecution?.status?.status === "ERROR" && (
                        <Alert variant="danger">{selectedExecution?.status?.error?.replace(/(?:\\n|\\\\n)/g, '\n')}</Alert>
                      )}
                      {selectedExecution?.status?.plugins && selectedExecution?.status?.plugins?.length && (
                        <SharedInfraViewDiagram
                          initialNodes={selectedExecution?.status?.plugins ? toNodes(selectedExecution?.status?.plugins) : []}
                          initialEdges={selectedExecution?.status?.plugins ? toEdges(selectedExecution?.status?.plugins) : []}
                        />
                      )}
                      
                    </Accordion.Body>
                  </Accordion.Item>
                ))}
              </Accordion>
            </div>
          </Tab>
        </Tabs>
      </Container>
    </>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default SharedInfraView