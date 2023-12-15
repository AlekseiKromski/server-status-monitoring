import styles from './main.module.css';
import {useEffect, useState} from "react";
import * as React from "react";
import TableContainer from "@mui/material/TableContainer";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import TableBody from "@mui/material/TableBody";
import Paper from "@mui/material/Paper";
import {Chip, styled, TableCell, tableCellClasses} from "@mui/material";
import moment from "moment"

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

function lastCheck(last) {
    return moment().diff(moment(last), 'second')
}

export default function Main() {

    let [status, setStatus] = useState("unknow")
    let [applications, setApplications] = useState([])

    useEffect(() => {
        let socket = new WebSocket(process.env.REACT_APP_SERVER);
        socket.onopen = () => {
            setStatus("connected")
        }

        socket.onmessage = function(event) {
            let applications = JSON.parse(event.data)
            setApplications(applications)
        };

        socket.onclose = function(event) {
            setStatus("connection closed")
        };

        socket.onerror = function(error) {
            setStatus("error: " + error)
        };
    }, [])

    return (
        <div className={styles.main}>
            <div className={styles.main}>
                <div className={styles.status}>
                    <span>Status: <Chip label={status} color="primary" /></span>
                </div>
                <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 1000 }} aria-label="customized table">
                        <TableHead>
                            <TableRow>
                                <StyledTableCell>Application name</StyledTableCell>
                                <StyledTableCell align="right">Url</StyledTableCell>
                                <StyledTableCell align="right">Status</StyledTableCell>
                                <StyledTableCell align="right">Last check</StyledTableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {
                                applications.map(app => {
                                    let appParsed = JSON.parse(app)
                                    return(
                                        <StyledTableRow>
                                            <StyledTableCell component="th" scope="row"></StyledTableCell>
                                            <StyledTableCell align="right">
                                                <a href={appParsed.link}>{appParsed.link}</a>
                                            </StyledTableCell>
                                            <StyledTableCell align="right">
                                                {appParsed.status === "SUCCESS" ? <Chip label="success" color="success" /> : <Chip label="error" color="error" />}
                                            </StyledTableCell>
                                            <StyledTableCell align="right">
                                                {lastCheck(appParsed.lastCheck)} seconds
                                            </StyledTableCell>
                                        </StyledTableRow>
                                    )
                                })
                            }
                        </TableBody>
                    </Table>
                </TableContainer>
            </div>
        </div>
    );
}
