import React, { memo, useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import "./index.css"
import InfraDiagram from "../InfraDiagram";
import { toEdges, toNodes } from "../InfraDiagram/utils";
import DefaultPanel from "./DefaultPanel";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";
import ReactAce from "react-ace/lib/ace";

const getBadgeVariants = (status: string) => {
  if (status === "RUNNING") {
    return 'primary'
  }

  if (status === "SUCCESS") {
    return 'success'
  }


  return 'danger'
}

const InfraEditor = memo(() => {
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()
  const { name } = useParams()
  const [infra, setInfra] = useState()
  const [currentInfra, setCurrentInfra] = useState<any>({ name: '', namespace: 'default', providerConfigRef: '' })
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [internalNodes, setInternalNodes] = useState([]);
  const [internalEgdes, setInternalEdges] = useState([]);
  const [tasks, setTasks] = useState<any>([])

  const getInfra = useCallback(async (name: string) => {
    const infraRes = await fetch(`http://localhost:8080/infra/${name}`)
    const infra = await infraRes.json()

    setInfra(infra)
    setNodes(toNodes(infra.tasks, "default"))
    setEdges(toEdges(infra.tasks, false))
  }, [])

  useEffect(() => {
    if (!name)
      return

    getInfra(name)
  }, [name])

  const createInfra = useCallback(async (infra: any) => {
    const res = await fetch(`http://localhost:8080/infra`, {
      method: 'POST',
      body: JSON.stringify({
        ...infra,
        tasks,
      })
    })
    const created = await res.json()
    navigate(`/infra/${infra?.name}`)
  }, [tasks])

  const handleDiagramChanges = useCallback((nodes: any, edges: any) => {
    let dict: any = {}
    for(let i = 0; i < nodes.length; i++) {
      dict[nodes[i].id] = nodes[i]?.data?.name
    }

    const newTasks = nodes?.map((node: any) => {
      return {
        name: node?.data?.name,
        ref: node?.data?.ref,
        type: node?.data?.type,
        depends: edges.filter((e: any) => e.target === node.id).map((e: any) => dict[e.source]),
        inputs: node?.data?.inputs,
        outputs: [],
      }
    })

    setTasks(newTasks)
    setCurrentInfra((s: any) => ({ ...s, tasks: newTasks }))

  }, [setTasks])

  const handleCodeChanges = useCallback((rawInfra: any) => {
   try {
    const parsed = JSON.parse(rawInfra)
    setCurrentInfra(parsed)
   } catch {

   }

  }, [])

  useEffect(() => {
    handleDiagramChanges(nodes, edges)
  }, [nodes, edges])

  useEffect(() => {
    setCurrentInfra(infra)
  }, [infra])

  useEffect(() => {
    console.log('CHANGE')
    setNodes(internalNodes)
    setEdges(internalEgdes)
  }, [searchParams])

  useEffect(() => {

    setInternalNodes(toNodes(currentInfra?.tasks || [], "default"))
    setInternalEdges(toEdges(currentInfra?.tasks || [], false))
  }, [currentInfra])
  
  return (
    <div className="shared-infra-create__content">
      <DefaultPanel
        infra={infra}
        onChange={(e: any) => setCurrentInfra({...e, tasks})}
        onSave={createInfra}
        goToView={() => navigate(`/infra/${name}`)}
      />
      <div className="shared-infra-view__diagram">
      {searchParams.get("view") === 'DIAGRAM' && (
        <InfraDiagram
          action="CREATE"
          nodes={nodes}
          edges={edges}
          onChangeDiagram={handleDiagramChanges}
        />
      )}
      {searchParams.get("view") === 'CODE' && (
        <ReactAce
          mode="json"
          theme="monokai"
          width='100%'
          height='200px'
          value={JSON.stringify(currentInfra, null, 2)}
          onChange={handleCodeChanges}
          tabSize={2}
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
})

export default InfraEditor