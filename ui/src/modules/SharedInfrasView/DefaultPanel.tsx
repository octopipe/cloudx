import React from 'react'
import './index.css'
import { Alert, Card, ListGroup, ListGroupItem, Nav, Spinner } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const getClassNameByExecution = (execution: any) => {
  if (execution?.status === "SUCCESS") {
    return 'shared-infra-diagram__default-panel__execution--success'
  }

  if (execution?.status === "RUNNING") {
    return 'shared-infra-diagram__default-panel__execution--running'
  }

  if (execution?.status === "ERROR") {
    return 'shared-infra-diagram__default-panel__execution--error'
  }

  return 'shared-infra-diagram__default-panel__execution--running'
}

const DefaultPanel = ({ sharedInfra, executions, onViewClick, onEditClick, onReconcileClick, onSelectExecution, onClose }: any) => {
  return (
    <div className='shared-infra-diagram__default-panel'>
      <div>
        <div className='d-flex justify-content-between'>
          <Card.Title>{sharedInfra?.name}</Card.Title>
          <div>
            <FontAwesomeIcon style={{cursor: 'pointer'}} icon="diagram-project" onClick={onViewClick} />
            <FontAwesomeIcon style={{cursor: 'pointer'}} className='ms-2' icon="edit" onClick={onEditClick} />
            <FontAwesomeIcon style={{cursor: 'pointer'}} className='ms-2' icon="rotate" onClick={onReconcileClick} />
          </div>
        </div>
        <p>{sharedInfra?.description}</p>
        <div className='mb-3'>
          <strong>Provider Config</strong><br/>
          {sharedInfra?.providerConfigRef?.name}
        </div>
        <div className='mb-3'>
          <strong>Plugins</strong><br/>
          {sharedInfra?.plugins?.map((p: any) => (
            <Card className='p-2 mt-2'>
              {p?.name}
            </Card>
          ))}
        </div>
        <div>
          <strong>Last execution</strong><br/>
          <small className='mb-2'>If you want to see more executions try to use webhooks to listen and save these events</small>
          {/* {executions?.map((i: any) => ( */}
            <Card 
              className={`${getClassNameByExecution(sharedInfra?.status)} mt-2`}
              style={{cursor: 'pointer'}} 
              onClick={() => onSelectExecution(sharedInfra?.status)}
            >
              <div className='d-flex'>
              {sharedInfra?.status?.status === "RUNNING" && <Spinner size='sm' className='me-1' />}
              {sharedInfra?.name}
              </div>
              
            </Card>
          {/* ))} */}
          
        </div>
      </div>
    </div>
  )
  
}

export default DefaultPanel