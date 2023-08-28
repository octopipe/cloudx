import React, { useCallback, useEffect, useState } from "react"
import { Background, ReactFlow, useEdgesState, ReactFlowProvider, useNodesState, useReactFlow } from "reactflow"
import { useFetch } from "use-http";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { NavLink, Outlet, useLocation, useNavigate, useParams, useSearchParams } from "react-router-dom";
import AceEditor from "react-ace";
import { Alert, Button, Card, Col, Container, Form, Modal, Nav, Navbar, Row, Tab, Tabs } from "react-bootstrap";
import 'reactflow/dist/style.css';
import ExecutionNode from "./ExecutionNode";
import DefaultNode from "./DefaultNode";

import { getLayoutedElements, toEdges, toNodes } from "./utils";
import './style.scss'
import { Ace } from "ace-builds";

const InfraView = ({ mode }: any) => {
  const location = useLocation()
  const { workspaceId, infraId } = useParams()
  const { fitView } = useReactFlow();
  const navigate = useNavigate()
  const [lastExecution, setLastExecution] = useState<any>()
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const { get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [infra, setInfra] = useState<any>()
  const [infraInterval, setInfraInterval] = useState<any>()
  const [infraEditorValue, setInfraEditorValue] = useState<any>('')
  const [hasModifications, setHasModifications] = useState<boolean>(false)
  const [show, setShow] = useState(false);
  const [currTask, setCurrTask] = useState<any>(null)

  useEffect(() => {
    getInfra()
    
    const interval = setInterval(() => {
      getInfra()
    }, 5000)

    setInfraInterval(interval)

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (infra && !lastExecution && mode == 'last-execution') {
      setLastExecution(infra?.status)
      return
    }

    setLastExecution(null)
  }, [infra, mode])

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
      <Navbar style={{height: '60px', borderBottom: '1px solid #ccc'}} className="text-white" bg="light">
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto ms-3">
            <div onClick={() => navigate(-1)} style={{cursor: 'pointer'}} className="nav-link nav-link-sub py-3 me-4 text-decoration-none">
              <FontAwesomeIcon icon="arrow-left" />
            </div>
            
            {/* <Nav.Link href="#home">Edit</Nav.Link>
            <Nav.Link href="#link">Link</Nav.Link> */}
          </Nav>
        </Navbar.Collapse>
      </Navbar>
      <div style={{ height: 'calc(100vh - 60px)', display: 'flex' }}>
        
        <div
          style={{ width: '40%',  maxHeight: 'calc(100vh - 60px)', zIndex: 10 }}
          className="p-2 bg-light"
        >
          <Nav variant="pills" defaultActiveKey={mode}>
            <Nav.Item>
              <Nav.Link eventKey="info" onClick={() => navigate(`/workspaces/${workspaceId}/infras/${infra?.name}`)}>Info</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="last-execution" onClick={() => navigate(`/workspaces/${workspaceId}/infras/${infra?.name}/last-execution`)}>Last execution</Nav.Link>
            </Nav.Item>
          </Nav>
          <div className="mt-2">
            {mode === 'info' && (
              <AceEditor
                style={{ width: '100%', height: 'calc(100vh - 125px)' }}
                mode="json"
                theme="github"
                name="infra-editor"
                onChange={(value) => JSON.stringify(infra, null, 2)}
                fontSize={14}
                showPrintMargin={true}
                showGutter={true}
                highlightActiveLine={true}
                value={JSON.stringify(infra, null, 2)}
                setOptions={{
                  useWorker: false,
                  showLineNumbers: true,
                  tabSize: 2,
                }}
              />
            )}
            {mode === 'last-execution' && (
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
                )}
                {infra?.status?.status === "" && (
                  <Alert style={{ cursor: 'pointer' }} variant="secondary" className="my-2" onClick={() => setLastExecution(infra?.status)}>
                    Not executed yet

                  </Alert>
                ) }
                {infra?.status?.status === "SUCCESS" && (
                  <>
                    <strong>Last execution: </strong>
                    <Alert variant="success" className="my-2" onClick={() => setLastExecution(infra?.status)}>
                      The last execution executed successfully
                    </Alert>
                  </>
                )}
                {infra?.status?.status === "RUNNING" && (
                  <>
                    <strong>Last execution: </strong>
                    <Alert variant="primary" className="my-2" onClick={() => setLastExecution(infra?.status)}>
                      Running...
                    </Alert>
                  </>
                )}
              </>
            )}
          </div>
        </div>

        
        
        <div style={{ width: '60%', height: 'calc(100vh - 60px)' }}>  
          <ReactFlow
            nodes={nodes}
            edges={edges} 
            nodeTypes={{
              executionNode: ExecutionNode,
              defaultNode: DefaultNode,
            }}
            fitView
            fitViewOptions={{maxZoom: 1}}
            onNodeClick={(event, node) => navigate(`/workspaces/${workspaceId}/infras/${infra?.name}/last-execution/task/${node?.data?.name}`)}
          >
            <Background/>
          </ReactFlow>
        </div>
      </div>

      <Outlet context={[infra, setInfra]}  />
    </>
  )
}

export default (props: any) => (
  <ReactFlowProvider>
    <InfraView {...props} />
  </ReactFlowProvider>
)