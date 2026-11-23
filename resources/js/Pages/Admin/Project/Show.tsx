import { PortfolioShow, type Project } from '../Work/Portfolio'
export default function Show({ item }: { item: Project }) { return <PortfolioShow kind="Project" item={item} /> }
