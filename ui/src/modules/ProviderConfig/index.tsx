import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Container, Form, ListGroup, Modal } from "react-bootstrap";
import AceEditor from "react-ace";
import { Link } from "react-router-dom";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";

const initialValue = {
  "name": "aws-config",
  "namespace": "default",
  "type": "AWS",
  "config": {
    "source": "CurrentAccount"
  }
}

const ProvidersConfigs = () => {
  const [list, setList] = useState<any>()
  const [selectedNode, setSelectedNode] = useState<any>()
  const [action, setAction] = useState<string>("VIEW")
  const [ value, setValue ] = useState<string>("")
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);

  const getList = useCallback(async () => {
    const res = await fetch("http://localhost:8080/providers-configs")
    const list = await res.json()

    setList(list)
  }, [])

  const saveChanges = useCallback(async () => {
    const method = action == "CREATE" ? "POST" : "PUT"
    const path = action == "CREATE" ? "" : `/${selectedNode.name}`
    const res = await fetch(`http://localhost:8080/providers-configs${path}`, {
      method,
      body: value, 
    })

    getList()
    handleClose()
  }, [value, action, selectedNode])

  const deleteNode = useCallback(async () => {
    const path = `/${selectedNode.name}`
    const res = await fetch(`http://localhost:8080/providers-configs${path}`, {
      method: "DELETE",
    })

    getList()
    handleClose()
  }, [value, action, selectedNode])

  

  useEffect(() => {
    getList()
  }, [])
  
  return (
    <>
      <Container fluid>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Providers config</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <Button
              variant="secondary"
              onClick={() => {
                setAction("CREATE")
                setValue(JSON.stringify(initialValue, null, 2))
                handleShow()
              }}
            >
              Create
            </Button>
          </div>
        </div>
        <ListGroup>
          {list?.items?.map((item: any, idx: any) => (
          <ListGroup.Item
            action
            className="d-flex justify-content-between align-items-start"
            onClick={() => {
              setSelectedNode(item)
              setValue(JSON.stringify(selectedNode, null, 2))
              setAction("VIEW")
              handleShow()
            }}
          >
            <div className="ms-2 me-auto">
              <div className="fw-bold">
                {item?.name}
              </div>
              {item?.description}
            </div>
            <Badge bg="primary" pill>
              {item?.spec?.plugins?.length}
            </Badge>
          </ListGroup.Item>
          ))}
        </ListGroup>
      </Container>
      <Modal show={show} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>
            {action === "CREATE" ? "New connection interface" : selectedNode?.name}
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <div className="d-flex justify-content-end mb-3">
            {action == "VIEW" && (
              <span
                style={{cursor: "pointer"}}
                onClick={() => setAction("EDIT")}
              >
                Edit
              </span>
            )}
            {action == "EDIT" && (
              <span
                style={{cursor: "pointer"}}
                onClick={() => setAction("VIEW")}
              >
                View
              </span>
            )}

            <span
              style={{cursor: "pointer"}}
              onClick={deleteNode}
              className="ms-2"
            >
              Delete
            </span>
            
          </div>
          {action == "VIEW" && (
            <>
              <div><strong>Type: </strong>{selectedNode?.type}</div>
              <div><strong>Source: </strong>{selectedNode?.config?.source}</div>              
            </>
          )}

          {(action == "EDIT" || action == "CREATE") && (
            <AceEditor
              mode="json"
              theme="github"
              width="100%"
              height="300px"
              onChange={value => setValue(value)}
              value={value}
              editorProps={{ $blockScrolling: false }}
            />
          )}
          
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Close
          </Button>
          {(action == "EDIT" || action == "CREATE") && (
            <Button variant="primary" onClick={saveChanges}>
              Save Changes
            </Button>
          )}
          
        </Modal.Footer>
      </Modal>
    </>
  )
}

export default ProvidersConfigs