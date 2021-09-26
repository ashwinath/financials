import React from 'react';

import {
  EuiSideNav,
  EuiIcon,
} from '@elastic/eui';

import { useLocation } from 'react-router-dom';
import { useHistory } from "react-router-dom";

const HOME_PAGE = "Home";

const PATH_MAPPING = {
  [HOME_PAGE]: "/",
}

export default function SideBar() {
  const history = useHistory();
  const { pathname } = useLocation();

  const sideNav = [
    {
      name: 'Financials',
      icon: <EuiIcon type="logoElasticsearch" />,
      id: '0',
      items: [
        {
          name: HOME_PAGE,
          id: '0.1',
          onClick: () => history.push(PATH_MAPPING[HOME_PAGE]),
          isSelected: pathname === PATH_MAPPING[HOME_PAGE],
        },
      ],
    },
  ];

  return (
		<EuiSideNav
			mobileTitle="Nav Items"
			items={sideNav}
		/>
  );
}
