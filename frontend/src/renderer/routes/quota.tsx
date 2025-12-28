import { createFileRoute } from '@tanstack/react-router'
import Quota from '../pages/Quota'

export const Route = createFileRoute('/quota')({
  component: Quota,
})

