import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Provider } from 'use-http';
import { ToastContainer, toast } from 'react-toastify';
import Workspaces from './modules/Workspaces';
import WorkspaceSettings from './modules/WorkspaceSettings';
import Workspace from './modules/Workspace';
import Infras from './modules/Infras';
import InfraView from './modules/InfraView';
import TasksOutputs from './modules/TasksOutputs';
import ProvidersConfig from './modules/ProvidersConfig';
import Webhooks from './modules/Webhooks';
import ProvidersConfigView from './modules/ProviderConfigView';
import { Alert } from 'react-bootstrap';

const RequestErrorToast = ({ res }: any) => (
  <div>
    <div><strong>url</strong>: {res?.url}</div>
    <div><strong>status</strong>: {res?.status}</div>
    <div><strong>data</strong>: {JSON.stringify(res?.data, null, 2)}</div>

  </div>
)


const App = () => {
  const [errors, setErrors] = React.useState<any>({})

  const options = {
    interceptors: {
      request: async ({ options, url, path, route }: any) => {
        return options
      },
      response: async ({ response }: any) => {
        const res = response
        if (!res.ok) {
          setErrors((errors: any) => ({ ...errors, [res.url]: res }))
          // toast(<RequestErrorToast res={res} />, {
          //   position: toast.POSITION.TOP_RIGHT,
          //   theme: 'colored',
          //   type: 'error',
          //   pauseOnFocusLoss: false,
          //   hideProgressBar: true,
          // })
        }
        return response
      }
    }
  }

  return (
    <Provider url="http://localhost:8080" options={options}>
      {errors.length > 0 && Object.keys(errors)?.map((url: any) =>
        <Alert className='m-0' variant='danger'>
          <div>
            <div><strong>url</strong>: {errors[url]?.url}</div>
            <div><strong>status</strong>: {errors[url]?.status}</div>
            <div><strong>data</strong>: {JSON.stringify(errors[url]?.data, null, 2)}</div>
          </div>
        </Alert>
      )}
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Workspaces />} />
          <Route path="/workspaces/:workspaceId" element={<Workspace />}>
            <Route path='' index element={<Navigate to="infras" />} />
            <Route path='infras' element={<Infras />} />
            <Route path='infras/create' element={<InfraView />} />
            <Route path='infras/:infraId' element={<InfraView />} />
            <Route path='tasks-outputs' element={<TasksOutputs />} />
            <Route path='providers-config' element={<ProvidersConfig />} />
            <Route path='providers-config/:providerConfigId' element={<ProvidersConfigView />} />
            <Route path='webhooks' element={<Webhooks />} />
            <Route path="settings" element={<WorkspaceSettings />} />
          </Route>
        </Routes>
      </BrowserRouter>
      <ToastContainer />
    </Provider>
  );
}

export default App
