import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React, { useCallback, useEffect, useState } from "react"
import { Alert, Breadcrumb, Button, ListGroup, ProgressBar, Spinner } from "react-bootstrap"
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { useFetch } from "use-http"


const Infras = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { workspaceId } = useParams()
  const { patch, get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [infras, setInfras] = useState<any>([])
  const [timer, setTimer] = useState<any>(null)

  useEffect(() => {
    getInfras()
    const interval = setInterval(() => {
      getInfras()
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const getProgressNow = (infra: any) => {
    return  (infra?.status?.tasks?.length * 100) / infra?.tasks?.length
  }
  
  const getInfras = useCallback(async () => {
    const infras = await get(`/infra`)
    if (response.ok) setInfras(infras?.items)
  }, [get])

  const sync = useCallback(async (infraId: string) => {
    const infras = await patch(`/infra/${infraId}/reconcile`)
  }, [get])

  return (
    <div className="m-4">
      <div>
        <Breadcrumb>
          {location.pathname.split('/').map((path, index) => {
            if (path === '') {
              return null
            }

            return (
              <Breadcrumb.Item key={index} href="#">
                {path}
              </Breadcrumb.Item>
            )
          })}
        </Breadcrumb>
      </div>
      <div className="mb-3 d-flex align-items-center justify-content-between">
        <h1>Infras</h1>
        <Button onClick={() => navigate(`/workspaces/${workspaceId}/infras/create`)}>
          <FontAwesomeIcon icon="plus" /> Create
        </Button>
      </div>
      <ListGroup variant="flush">
        {infras.map((infra: any, index: any) => (
          <ListGroup.Item className="mb-2">
            <div className="ms-2 me-auto">

              <div className="d-flex justify-content-between">
                <div>
                  <strong>Name: </strong>
                  <a 
                    className="fw-bold"
                    style={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/workspaces/${workspaceId}/infras/${infra?.name}?mode=VIEW`)}
                  >
                    {infra?.name}
                  </a>
                </div>
                <div>
                {infra?.status?.status === "RUNNING" ? (
                  <Spinner variant="primary" animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                  </Spinner>
                ) : (
                  <>
                    <FontAwesomeIcon className="me-2" icon="rotate" style={{ cursor: 'pointer' }} onClick={() => sync(infra?.name)} />
                    <FontAwesomeIcon icon="trash" color="red" style={{ cursor: 'pointer' }} />
                  </>
                )}
                  
                </div>
              </div>
              <strong>Description: </strong>{infra?.description}<br/><br/>
              
              
              {infra?.status?.status === "ERROR" && (
                <>
                  <strong>Last execution: </strong>
                  <Alert variant="danger" className="my-2" onClick={e => e.preventDefault()}>
                    <strong>Error: </strong>{infra?.status?.error?.message}<br/>
                    <strong>Code: </strong>{infra?.status?.error?.code}<br/>
                    <strong>Message: </strong>{infra?.status?.error?.tip}<br/>
                  </Alert>
                </>
              )}
              {infra?.status?.status === "SUCCESS" && (
                <>
                  <strong>Last execution: </strong>
                  <Alert variant="success" className="my-2" onClick={e => e.preventDefault()}>
                    The last execution executed successfully
                  </Alert>
                </>
              )}
            </div>
            

          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  )
}

export default Infras