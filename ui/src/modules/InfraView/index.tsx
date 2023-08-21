import React, { useCallback, useEffect, useState } from "react"
import { Background, ReactFlow, useEdgesState, useNodesState } from "reactflow"
import { useFetch } from "use-http";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { NavLink, useLocation, useNavigate, useParams, useSearchParams } from "react-router-dom";
import AceEditor from "react-ace";
import { Alert, Button, Card, Col, Container, Form, Nav, Navbar, Row, Tab, Tabs } from "react-bootstrap";
import 'reactflow/dist/style.css';
import ExecutionNode from "./ExecutionNode";
import DefaultNode from "./DefaultNode";

import { getLayoutedElements, toEdges, toNodes } from "./utils";
import './style.scss'

const InfraView = () => {
  const location = useLocation()
  const { workspaceId, infraId } = useParams()
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()
  const [lastExecution, setLastExecution] = useState<any>()
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const { get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [infra, setInfra] = useState<any>()
  const [infraInterval, setInfraInterval] = useState<any>()
  const [infraEditorValue, setInfraEditorValue] = useState<any>('')
  const [hasModifications, setHasModifications] = useState<boolean>(false)

  useEffect(() => {
    getInfra()
    if (searchParams.get('mode') === 'EDIT') {
      clearInterval(infraInterval)
      return
    }
    
    const interval = setInterval(() => {
      getInfra()
    }, 5000)

    setInfraInterval(interval)

    return () => clearInterval(interval)
  }, [searchParams])

  const setNodesAndEdges = useCallback((n: any[], e: any[], nodeType: string, animated: boolean) => {
    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
      toNodes(n, nodeType),
      toEdges(e, animated),
    );

    setNodes(layoutedNodes);
    setEdges(layoutedEdges);
  }, [])

  useEffect(() => {
    try {
      console.log('EDIT')
      const infraEditorObject = JSON.parse(infraEditorValue)
      const newInfra = { ...infra, ...infraEditorObject }
      setNodesAndEdges(newInfra?.tasks, newInfra?.tasks, "defaultNode", false)
      setHasModifications(true)
    } catch(error) {
      console.log(error)
    }
  }, [infraEditorValue])

  useEffect(() => {
    if (!infra) return

    if (!!lastExecution) {
      setLastExecution(infra?.status)
      return
    }

    setInfraEditorValue(JSON.stringify({
      description: infra?.description,
      providerConfigRef: infra?.providerConfigRef,
      runnerConfigRef: infra?.runnerConfigRef,
      tasks: infra?.tasks,
    }, null, 2))
    setNodesAndEdges(infra?.tasks, infra?.tasks, "defaultNode", false)
  }, [infra])

  useEffect(() => {
    if (infra && !lastExecution) {
      setNodesAndEdges(infra?.tasks, infra?.tasks, "defaultNode", false)
      return
    }

    if (lastExecution)
      setNodesAndEdges(lastExecution?.tasks, lastExecution?.tasks, "executionNode", true)
  }, [lastExecution])

  const getInfra = useCallback(async () => {
    const infra = await get(`/infra/${infraId}`)
    if (response.ok) setInfra(infra)
  }, [get])

  return (
    <>
      <Navbar style={{height: '60px'}} className="text-white" bg="dark" data-bs-theme="dark">
        <Container>
          
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="me-auto">
              <NavLink to={`/workspaces/${workspaceId}/infras`}  className="nav-link-sub py-3 me-4 text-decoration-none">
                Voltar
              </NavLink>
              <NavLink to={`/workspaces/${workspaceId}/infras/${infraId}?mode=EDIT`} className="nav-link-sub py-3 text-decoration-none">
                Editar
              </NavLink>
              {/* <Nav.Link href="#home">Edit</Nav.Link>
              <Nav.Link href="#link">Link</Nav.Link> */}
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
      <div style={{ height: '100vh', display: 'flex' }}>
        
        {/* <div className="d-flex flex-column justify-content-between text-white" style={{width: '4.5rem', borderRight: '1px solid #ccc'}}> */}
          {/* <Nav variant="pills" activeKey={location.pathname} className="nav-flush flex-column mb-auto text-center">
            <NavLink to={`/workspaces/${workspaceId}/infras`} className="nav-link-sub text-black py-3">
              <FontAwesomeIcon size="lg" icon="arrow-left" />
            </NavLink>
            <NavLink to={`/workspaces/${workspaceId}/infras/${infraId}?mode=${searchParams.get("mode")}`} className="nav-link-sub text-black py-3">
              <FontAwesomeIcon size="lg" icon="edit" />
            </NavLink>
          </Nav> */}
          
          {/* <FontAwesomeIcon
            icon={infra && !lastExecution ? "arrow-left" : "close"}
            onClick={() => infra && !lastExecution ? navigate(-1) : setLastExecution(null)}
            cursor={'pointer'}
            size={infra && !lastExecution ? "sm" : "2x"}
          /> */}
        {/* </div> */}
        {searchParams.get('mode') === 'EDIT' && (
          <div className="bg-light" style={{width: infra && !lastExecution ? '45%' : '5%', maxHeight: '100vh'}}>
            {infra && !lastExecution && (
              <>
                <Tabs className="mt-2" style={{width: '100%'}}>
                  <Tab style={{height: '100%'}} eventKey="code" title="Code">
                    <AceEditor
                      mode="json"
                      theme="github"
                      onChange={v => setInfraEditorValue(v)}
                      name="UNIQUE_ID_OF_DIV"
                      editorProps={{ $blockScrolling: true }}
                      setOptions={{
                        useWorker: false,
                        readOnly: (!searchParams.has('mode') || searchParams?.get('mode') === 'VIEW'),
                      }}
                      value={infraEditorValue}
                      width="100%"
                      height="87%"
                    />
                    <div style={{position: 'fixed', bottom: 0, zIndex: 5, width: '100%'}}>
                      <Button style={{borderRadius: 0}}>Save</Button>
                    </div>
                  </Tab>
                  <Tab className="p-4" eventKey="last-execution" title="Last execution">
                    <div className="bg-light">
                      {infra && !lastExecution && (
                        <>
                          <h5>Last execution</h5>
                          <small>Cloudx save only last execution, if you want to save more, use webhooks.</small>
                          {infra?.status?.status === "ERROR" && (
                            <Alert style={{ cursor: 'pointer' }} variant="danger" className="my-2" onClick={() => setLastExecution(infra?.status)}>
                              <strong>Error: </strong>{infra?.status?.error?.message}<br/>
                              <strong>Code: </strong>{infra?.status?.error?.code}<br/>
                              <strong>Message: </strong>{infra?.status?.error?.tip}<br/>
                              <strong>Started At: </strong>{infra?.status?.startedAt}<br/>

                            </Alert>
                          ) }
                        </>
                      )}
                      {lastExecution && (
                        <div style={{width: '100%'}}>
                          <Button variant="decondary" onClick={() => setLastExecution(null)}>Exit execution mode</Button>
                        </div>
                      )}
                    </div>
                  </Tab>
                </Tabs>
              </>
            )}
            
          </div>
        )}
        
        <div style={{ width: '55%', height: '100vh' }}>
          <ReactFlow
            nodes={nodes}
            edges={edges} 
            nodeTypes={{
              executionNode: ExecutionNode,
              defaultNode: DefaultNode,
            }}
            fitView
            fitViewOptions={{maxZoom: 1}}
          >
            <Background/>
          </ReactFlow>
        </div>
      </div>
    </>
  )
}

export default InfraView