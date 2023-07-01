import React from 'react'
import './index.css'
import { Card, ListGroup, ListGroupItem, Nav } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const DefaultPanel = ({ sharedInfra, onSelectExecution, onClose }: any) => {
  const onDragStart = (event: any, nodeType: any) => {
    event.dataTransfer.setData('application/reactflow', nodeType);
    event.dataTransfer.effectAllowed = 'move';
  };

  return (
    <Card className='shared-infra-diagram__default-panel'>
      <Card.Body>
        <Card.Title>{sharedInfra?.name}</Card.Title>
        <p>{sharedInfra?.description}</p>
        <div>
          {sharedInfra?.status?.executions?.map((i: any) => (
            <Card className='mt-2' style={{cursor: 'pointer'}} onClick={() => onSelectExecution(i)}>
              <Card.Body>{i.name}</Card.Body>
            </Card>
          ))}
          
        </div>
      </Card.Body>
    </Card>
  )
  
}

export default DefaultPanel