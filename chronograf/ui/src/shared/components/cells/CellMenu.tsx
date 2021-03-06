// Libraries
import React, {PureComponent} from 'react'
import classnames from 'classnames'

// Components
import MenuTooltipButton, {
  MenuItem,
} from 'src/shared/components/MenuTooltipButton'

// Types
import {Cell, CellQuery} from 'src/types/v2/dashboards'

import {ErrorHandling} from 'src/shared/decorators/errors'

interface Props {
  cell: Cell
  isEditable: boolean
  dataExists: boolean
  onEdit: () => void
  onClone: (cell: Cell) => void
  onDelete: (cell: Cell) => void
  onCSVDownload: () => void
  queries: CellQuery[]
}

interface State {
  subMenuIsOpen: boolean
}

@ErrorHandling
class CellMenu extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props)

    this.state = {
      subMenuIsOpen: false,
    }
  }

  public render() {
    return <div className={this.contextMenuClassname}>{this.renderMenu}</div>
  }

  private get renderMenu(): JSX.Element {
    const {isEditable} = this.props

    if (isEditable) {
      return (
        <div className="dash-graph-context--buttons">
          {this.pencilMenu}
          <MenuTooltipButton
            icon="duplicate"
            menuItems={this.cloneMenuItems}
            informParent={this.handleToggleSubMenu}
          />
          <MenuTooltipButton
            icon="trash"
            theme="danger"
            menuItems={this.deleteMenuItems}
            informParent={this.handleToggleSubMenu}
          />
        </div>
      )
    }
  }

  private get pencilMenu(): JSX.Element {
    const {queries} = this.props

    if (!queries.length) {
      return
    }

    return (
      <MenuTooltipButton
        icon="pencil"
        menuItems={this.editMenuItems}
        informParent={this.handleToggleSubMenu}
      />
    )
  }

  private get contextMenuClassname(): string {
    const {subMenuIsOpen} = this.state

    return classnames('dash-graph-context', {
      'dash-graph-context__open': subMenuIsOpen,
    })
  }

  private get editMenuItems(): MenuItem[] {
    const {dataExists, onCSVDownload} = this.props

    return [
      {
        text: 'Configure',
        action: this.handleEditCell,
        disabled: false,
      },
      {
        text: 'Download CSV',
        action: onCSVDownload,
        disabled: !dataExists,
      },
    ]
  }

  private get cloneMenuItems(): MenuItem[] {
    return [{text: 'Clone Cell', action: this.handleCloneCell, disabled: false}]
  }

  private get deleteMenuItems(): MenuItem[] {
    return [{text: 'Confirm', action: this.handleDeleteCell, disabled: false}]
  }

  private handleEditCell = (): void => {
    const {onEdit} = this.props
    onEdit()
  }

  private handleDeleteCell = (): void => {
    const {onDelete, cell} = this.props
    onDelete(cell)
  }

  private handleCloneCell = (): void => {
    const {onClone, cell} = this.props
    onClone(cell)
  }

  private handleToggleSubMenu = (): void => {
    this.setState({subMenuIsOpen: !this.state.subMenuIsOpen})
  }
}

export default CellMenu
