import React from 'react'
import './index.css'
import { Alert, Card, ListGroup, ListGroupItem, Nav } from 'react-bootstrap'
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

const DefaultPanel = ({ sharedInfra, executions, onViewClick, onEditClick, onSelectExecution, onClose }: any) => {
  return (
    <div className='shared-infra-diagram__default-panel'>
      <div>
        <div className='d-flex justify-content-between'>
          <Card.Title>{sharedInfra?.name}</Card.Title>
          <div>
            <FontAwesomeIcon style={{cursor: 'pointer'}} icon="diagram-project" onClick={onViewClick} />
            <FontAwesomeIcon style={{cursor: 'pointer'}} className='ms-2' icon="edit" onClick={onEditClick} />
            <FontAwesomeIcon style={{cursor: 'pointer'}} className='ms-2' icon="rotate" />
          </div>
        </div>
        <p>{sharedInfra?.description}</p>
        <div>
          <strong>Executions</strong>
          {executions?.map((i: any) => (
            <Card 
              className={`${getClassNameByExecution(i)} mt-1`}
              style={{cursor: 'pointer'}} 
              onClick={() => onSelectExecution(i)}
            >
              {i.name}
            </Card>
          ))}
          
        </div>
      </div>
    </div>
  )
  
}

export default DefaultPanel