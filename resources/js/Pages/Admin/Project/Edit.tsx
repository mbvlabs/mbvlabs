import { PortfolioForm, type Project } from '../Work/Portfolio'
export default function Edit({ item }: { item: Project }) { return <PortfolioForm kind="Project" item={item} /> }
