import { PortfolioForm, type Work } from './Portfolio'
export default function Edit({ item }: { item: Work }) { return <PortfolioForm kind="Work" item={item} /> }
