import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React, { useCallback, useEffect, useState } from "react"
import { Breadcrumb, Button, Form, ListGroup } from "react-bootstrap"
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { useFetch } from "use-http"
import { useForm, SubmitHandler } from "react-hook-form"

const ProvidersConfigView = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { workspaceId, providerConfigId } = useParams()
  const { register, setValue, getValues, handleSubmit } = useForm<any>()
  const { patch, get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [providerConfig, setProviderConfig] = useState<any>([])


  const getProviderConfig = useCallback(async () => {
    const providersConfig = await get(`/providers-configs/${providerConfigId}`)
    if (response.ok) setProviderConfig(providersConfig)
  }, [get])

  useEffect(() => {
    getProviderConfig()
    const interval = setInterval(() => {
      getProviderConfig()
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    setValue('name', providerConfig?.name)
    setValue('source', providerConfig?.source)
  }, [providerConfig])

  return (
    <div className="p-4">
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
      <div style={{width: '40rem'}}>
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Name</Form.Label>
            <Form.Control {...register("name")} type="text"  placeholder="Infra name..." />
          </Form.Group>
          <div className="mb-3">
            <Form.Check
              inline
              label="AWS"
              name="type"
              type="radio"
              id="aws"
              checked={providerConfig?.type === 'AWS'}
            />
            <Form.Check
              inline
              label="Azure"
              name="group1"
              type="radio"
              id="azure"
              checked={providerConfig?.type === 'AZURE'}
            />
            <Form.Check
              inline
              label="GCP"
              name="group1"
              type="radio"
              id="gcp"
              checked={providerConfig?.type === 'GCP'}
            />
          </div>
          {providerConfig?.type === 'AWS' && (
            <>
              <Form.Group className="mb-3">
                <Form.Label>Source</Form.Label>
                <Form.Select {...register("source")} aria-label="Default select example">
                  <option>Open this select menu</option>
                  <option value="SECRET_REF">Secret Ref</option>
                  <option value="CREDENTIALS">Credentials</option>
                  <option value="ROLE_ARN">Role arn</option>
                </Form.Select>
              </Form.Group>
            </>
          )}
          {getValues()?.source === 'CREDENTIALS' && (
            <>
              <Form.Group className="mb-3">
                <Form.Label>Access key ID</Form.Label>
                <Form.Control type="text" value="***" placeholder="Infra name..." />
              </Form.Group>
              <Form.Group className="mb-3">
                <Form.Label>Secret access key</Form.Label>
                <Form.Control type="text" value="***" placeholder="Infra name..." />
              </Form.Group>
            </>
          )}
        </Form>
      </div>
    </div>
  )
}

export default ProvidersConfigView