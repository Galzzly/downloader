package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/colinmarc/hdfs/v2"
	"github.com/colinmarc/hdfs/v2/hadoopconf"
)

func GetAddresses(input string) ([]string, error) {
	f, err := ioutil.ReadFile(input)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(strings.TrimSpace(string(f)), "\n"), nil
}

func ConnectToNamenode() (*hdfs.Client, error) {
	namenode := os.Getenv("HADOOP_NAMENODE")
	conf, err := hadoopconf.LoadFromEnvironment()
	if err != nil {
		return nil, err
	}

	options := hdfs.ClientOptionsFromConf(conf)
	if namenode != "" {
		options.Addresses = []string{namenode}
	}

	if options.Addresses == nil {
		return nil, errors.New("cannot fine Namenode to connect to")
	}

	// This part of the code needs to be developed.
	// It will be commented out for the time being until I have
	// the time to review it and test it further.
	/*
		if options.KerberosClient != nil {
			options.KerberosClient, err = getKerberosClient() << function to be written
			if err != nil {
				return nil, fmt.Errorf("problem with kerberos auth: %s", err)
			}
		} else
	*/
	options.User = os.Getenv("HADOOP_USER_NAME")
	if options.User == "" {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("unable to determine user: %s", err)
		}
		options.User = u.Username
	}
	// }

	dialFunc := (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).DialContext

	options.NamenodeDialFunc = dialFunc
	options.DatanodeDialFunc = dialFunc

	client, err := hdfs.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to namenode: %s", err)
	}

	return client, nil
}

func DownloadFile(address string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	return client.Do(req)
}

func Complete() string {
	buf := new(bytes.Buffer)
	buf.WriteString(compmess)
	return buf.String()
}

const compmess = `&&@@@@@@@@@@@@@@@@@@@@@@@@@&*/ /( #@@@&@@@@&*#@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&
@@@@@@@@@@@@@@@@@@@@@@@@@/*%@@*,@@@%@/@(./##%&&@/@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@*                       (@@@@@@@@@@@@@@@&@&@@&...                      /@@@
@@@*    ////////////////@@@@@@@@@@@@@@@@@@@@@& *&, /&    ....,,,****//    /@@@
@@@*    @@@@@@@@@@@@@@&@@@@@@@@@@@@@@@@@@@@@@@@&&&@@@* @#@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@@&@@@@@@@@@@@@@@@@@@@@@@@@@@@*&&&#(* /@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@#*%###%@/ @@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@%(@@@@@@@@@@@@@@@@%#@@@@@@@@@@@*     .&.@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@#      .,.           .#@@@@@@@@*       # @@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@%       ./              %&#.  ../       #@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@      @@@(            &@@@@@@/.. ##   (@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@(/   /@@@@@@%.     /@@@@@@@@@,&&* @.#  @@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@%&@@@,@@@@@@@ %&@@@@@@@%,@@@@@@@@.,@@% %@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@,@(  &#/*,,/@#  . (@@@@@* @@@@@@@@..@ #@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@%@@@@% ,%@@@@@.  ##%@@@.(@@@@@/.,( @@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@@@@@@..%@@@@@@@@@@@@@@@@@@@.  * (@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@@          /@@@@@@@@@@@@@@ . ,(%@@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@/            ./@@@@@@@@@@@@&@.@@@@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@, @.  .@* &@@*,@& &@@@@&(@@...*%, @/@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@&  ,@@@@@&(* .     (&,@%#@&#&*. @@@&%@@@@@@@@@@@@@@@    /@@@
@@@*    (&&@@@@@@@@@ ,           *%@( &, #@@%.  %@@@@@#@@@@@@@@@@@@@@@    /@@@
@@@*    &@@@@@@/ */@@  .    &@@@@%% @ #*@. /@ %@@@@@@@@@@@%%@@@@@@@@@@    /@@@
@@@*    &/&     % &*@@.        .(.#.%&   %%.@@@@@@@@@@@@@@@@@@@@@@@@@@    /@@@
@@@*           /   #&&&..      %(.    .% %&&&&&&&&&&&%,%&&&&&&&&&&&%%,    /@@@
@@@*                                                                      /@@@
@@@*    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@/    ,@@@     ,@@@%     %@@      @@.     %@@@@@@@@@.    @@@@    /@@@
@@@*    @&  &@  *@@  .@*  @@%  @@@@@@  *@@@@@.  @#  &@@@@@@@%  @/ *@@@    /@@@
@@@*    @@*  &@@@@@  .@(  @@%  @@@@@@  *@@@@@.  @&  %@@@@@@@@  .  @@@@    /@@@
@@@*    @@@&  .@@@@      &@@%    %@@@     @@@.  @&  %@@@@@@@@(  (@@@@@    /@@@
@@@*    @@@@@.  @@@  .@@@@@@%  @@@@@@  *@@@@@.  @&  &@@@@@@@* ,  ,  %@    /@@@
@@@*    @/  @@  (@@  ,@@@@@@&  @@@@@@  *@@@@@,  @/  @@@@@@@@  %@/  .@@    /@@@
@@@*    @@,    *@@@. *@@@@@@&     %@@.     @@*    .%@@@@@@@@&    ,( *@    /@@@
@@@*    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    /@@@
@@@*    @/      @@@@%     (@@@. .@@  (@@  (@@.     .@@&      &@@@(  (@    /@@@
@@@*    @/  %@   @@&   @,  &@@/  @&   @#  @@@.  @@@@@@@  *@,  &@@&  &@    /@@@
@@@/    @/  &@.  &@&  *@*  &@@@  &(   @, .@@@.  @@@@@@@  /@#  %@@@  @@    /@@@
@@@/    @/  .   .@@%  *@*  %@@@  *    (  (@@@.     @@@@      .@@@@  @@    /@@@
@@@/    @/  .,%@@@@%  *@*  %@@@/    %    @@@@.  ##%@@@@   .  @@@@@  @@    /@@@
@@@/    @/  @@@@@@@%  ,@*  %@@@&   /@    @@@@.  @@@@@@@  /@  *@@@@ .@@    /@@@
@@@/    @/  @@@@@@@@   @   &@@@@   &@*  (@@@@.  #%%%@@@  /@,  @@@&%%@@    /@@@
@@@/    @/  &@@@@@@@@,   .@@@@@@/  @@@  &@@@@.     .@@@  /@&  /@@(  &@    /@@@
@@@/    %%%##############%######%###########%######%##%#####%###%%###%    /@@@
@@@/                                                                      /@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&`
