import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React, { useCallback, useEffect, useState } from "react"
import { Breadcrumb, Button, ListGroup } from "react-bootstrap"
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { useFetch } from "use-http"

const ProvidersConfig = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { workspaceId } = useParams()
  const { patch, get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [providersConfig, setprovidersConfig] = useState<any>([])

  const getProvidersConfig = useCallback(async () => {
    const providersConfig = await get(`/providers-configs`)
    if (response.ok) setprovidersConfig(providersConfig?.items)
  }, [get])

  useEffect(() => {
    getProvidersConfig()
    const interval = setInterval(() => {
      getProvidersConfig()
    }, 5000)

    return () => clearInterval(interval)
  }, [])

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
        <h1>Providers Config</h1>
        <Button><FontAwesomeIcon icon="plus" /> Create</Button>
      </div>
      <ListGroup variant="flush">
        {providersConfig.map((infra: any, index: any) => (
          <ListGroup.Item action onClick={() => navigate(`/workspaces/${workspaceId}/providers-config/${infra?.name}`)}>
            <div className="ms-2 me-auto">
              <div className="fw-bold">{infra?.name}</div>
              {infra?.description}
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  )
}

export default ProvidersConfig