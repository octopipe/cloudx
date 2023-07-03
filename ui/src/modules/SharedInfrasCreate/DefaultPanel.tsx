import React, { useCallback, useState } from 'react'
import './index.css'
import { Alert, Button, Card, Form, ListGroup, ListGroupItem, Nav } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const getClassNameByExecution = (execution: any) => {
  if (execution?.status?.status === "SUCCESS") {
    return 'shared-infra-diagram__default-panel__execution--success'
  }

  if (execution?.status?.status === "RUNNING") {
    return 'shared-infra-diagram__default-panel__execution--running'
  }

  return 'shared-infra-diagram__default-panel__execution'
}

const DefaultPanel = ({ sharedInfra, onCreate }: any) => {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  
  const handleCreate = useCallback(() => {
    onCreate({ name, namespace: "default", description })
  }, [name, description])

  return (
    <div className='shared-infra-diagram__default-panel'>
      <div>
        <Card.Title>{sharedInfra?.name}</Card.Title>
        <Form.Group className='mb-3'>
          <Form.Label>Name</Form.Label>
          <Form.Control type="text" placeholder='Type the of yout shared infra' value={name} onChange={e => setName(e.target.value)} />
        </Form.Group>
        <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1" >
          <Form.Label>Description</Form.Label>
          <Form.Control as="textarea" rows={3} value={description} onChange={e => setDescription(e.target.value) } />
        </Form.Group>
        <div className="d-grid gap-2">
          <Button onClick={handleCreate}>Create shared infra</Button>
        </div>
      </div>
    </div>
  )
  
}

export default DefaultPanel