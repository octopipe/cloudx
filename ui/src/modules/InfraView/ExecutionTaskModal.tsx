import React, { useCallback, useEffect, useState } from "react"
import { Background, ReactFlow, useEdgesState, ReactFlowProvider, useNodesState, useReactFlow } from "reactflow"
import { useFetch } from "use-http";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { NavLink, useLocation, useNavigate, useOutletContext, useParams, useSearchParams } from "react-router-dom";
import AceEditor from "react-ace";
import { Alert, Button, Card, Col, Container, Form, Modal, Nav, Navbar, Row, Tab, Tabs } from "react-bootstrap";
import 'reactflow/dist/style.css';
import ExecutionNode from "./ExecutionNode";
import DefaultNode from "./DefaultNode";

import { getLayoutedElements, toEdges, toNodes } from "./utils";
import './style.scss'
import { Ace } from "ace-builds";

const ExecutionTaskModal = () => {
  const navigate = useNavigate()

  const [infra] = useOutletContext<any>()
  const [infraInterval, setInfraInterval] = useState<any>()
  const { get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const { infraId, taskName } = useParams()
  const [currTask, setCurrTask] = useState<any>(null)

  useEffect(() => {
    setCurrTask(infra?.status?.tasks?.find((task: any) => task.name === taskName))
  }, [infra])

  return (
    <Modal show={true} onHide={() => ({})}>
      <Modal.Header closeButton>
        <Modal.Title>{currTask?.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <AceEditor
          style={{ width: '100%', height: '400px' }}
          mode="json"
          theme="github"
          name="infra-editor"
          fontSize={14}
          showPrintMargin={true}
          showGutter={true}
          highlightActiveLine={true}
          value={JSON.stringify(currTask, null, 2)}
          setOptions={{
            useWorker: false,
            showLineNumbers: true,
            tabSize: 2,
          }}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ExecutionTaskModal