import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import ReactAce from "react-ace/lib/ace";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import "./index.css"
import InfraDiagram from "../InfraDiagram";
import { toEdges, toNodes } from "../InfraDiagram/utils";
import DefaultPanel from "./DefaultPanel";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";


const getBadgeVariants = (status: string) => {
  if (status === "RUNNING") {
    return 'primary'
  }

  if (status === "SUCCESS") {
    return 'success'
  }


  return 'danger'
}

let interval: any

const InfraView = () => {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const { name: infraName } = useParams()
  const [infra, setInfra] = useState<any>()
  const [infraCode, setInfraCode] = useState('')
  const [executions, setExecutions] = useState<any>([])
  const [selectedExecution, setSelectedExecution] = useState<boolean>()
  const [heatmap, setHeatmap] = useState(false)
  const [nodes, setNodes] = useState<any>([])
  const [edges, setEdges] = useState<any>([])
  // const [codeView, setCodeView] = useState(false)


  const getInfra = useCallback(async (name: string) => {
    const infraRes = await fetch(`http://localhost:8080/infra/${name}`)
    const infra = await infraRes.json()

    setInfra(infra)
  }, [])
  
  const handleReconcile = useCallback(async () => {
    const res = await fetch(`http://localhost:8080/infra/${infraName}/reconcile`, {method: 'PATCH', body: JSON.stringify({})})
    const item = await res.json()
  }, [])

  useEffect(() => {
    setSearchParams({ view: 'DIAGRAM' })

    if (!infraName) return
    
    getInfra(infraName)
    interval = setInterval(() => {
      getInfra(infraName)
    }, 3000)
    

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    
    if (!infra) return

    const {status, ...rest} = infra
    setInfraCode(JSON.stringify(rest, null, 2))

    if (selectedExecution) {
      setNodes(toNodes(infra?.status?.tasks, "executionNode"))
      setEdges(toEdges(infra?.status?.tasks, true))
      return
    }

    setNodes(toNodes(infra?.tasks, "default"))
    setEdges(toEdges(infra?.tasks, false))
  }, [selectedExecution, infra])

  
  return (
    <div className="shared-infra-view__content">
      <DefaultPanel
        infra={infra}
        executions={executions}
        onViewCode={() => setSearchParams({ view: 'CODE' })}
        onViewClick={() => {
          setSearchParams({ view: 'DIAGRAM' })
          setSelectedExecution(false)
        }}
        onEditClick={() => navigate(`/infra/${infraName}/edit?view=${searchParams.get('view')}`)}
        onReconcileClick={() => handleReconcile()}
        onSelectExecution={(e: any) => {
          setSearchParams({ view: 'DIAGRAM' })
          setSelectedExecution(true)
        }}
      />
      {infra?.status && infra?.status?.error && (
        <Alert
          style={{position: 'fixed', top: '10px', right: '10px', left: '390px'}}
          variant="danger"
        >{infra?.status?.error}</Alert>
      )}
      {!!selectedExecution && (
        <div
          style={{
            position: "absolute",
            top: "20px",
            right: "20px",
            background: "#2c2c2e",
            color: "#fff",
            borderRadius: "50%",
            padding: "20px",
            width: "80px",
            height: "80px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 99,
          }}
        >
          <FontAwesomeIcon style={{cursor: 'pointer'}} size="2x" icon="fire" onClick={() => setHeatmap(!heatmap)} />
        </div>
      )}
      <div className="shared-infra-view__diagram">
      {searchParams.get("view") === "DIAGRAM" && <InfraDiagram
        infra={infra}
        nodes={nodes}
        edges={edges}
        isExecution={heatmap}
      />}
      {searchParams.get("view") === "CODE" && (
        <ReactAce
          mode="json"
          theme="monokai"
          width='100%'
          height='200px'
          value={infraCode}
          onChange={() => {}}
          tabSize={2}
          readOnly
          enableBasicAutocompletion={true}
          style={{
            width: '100vw',
            height: '100vh'
          }}
        />
      )}
      </div>
     
    </div>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default InfraView