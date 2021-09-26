import React from 'react';

import {
  EuiSideNav,
  EuiIcon,
} from '@elastic/eui';

import {
  useLocation,
  useHistory,
} from 'react-router-dom';

import { 
  useDispatch,
  useSelector,
} from 'react-redux';

import {
  logoutAsync,
} from '../redux/loginSlice';

const HOME_PAGE = "Home";

const PATH_MAPPING = {
  [HOME_PAGE]: "/",
}

export default function SideBar() {
  const history = useHistory();
  const dispatch = useDispatch();
  const { pathname } = useLocation();
  const { isLoggedIn } = useSelector((state) => state.login)

  if (!isLoggedIn) {
    history.push("/login");
  }

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
    {
      name: 'Account',
      icon: <EuiIcon type="logoElasticsearch" />,
      id: '0',
      items: [
        {
          name: "Log me out",
          id: '0.1',
          onClick: () => dispatch(logoutAsync()),
          isSelected: false,
        },
      ],
    }
  ];

  return (
		<EuiSideNav
			mobileTitle="Nav Items"
			items={sideNav}
		/>
  );
}
